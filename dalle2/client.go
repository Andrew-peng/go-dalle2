package dalle2

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

const (
	baseUrl  = "https://api.openai.com"
	goDalle2 = "go-dalle2"

	// headers
	headerAuth        = "Authorization"
	headerContentType = "Content-Type"
	headerUserAgent   = "User-Agent"

	applicationJson = "application/json"
	multipartForm   = "multipart/form-data"
)

type Client interface {
	Create(context.Context, string, ...Option) (*Response, error)
	Edit(context.Context, []byte, []byte, string, ...Option) (*Response, error)
	Variation(context.Context, []byte, ...Option) (*Response, error)
}

type ClientV1 struct {
	version string
	client  *http.Client

	Url       *url.URL
	UserAgent string
}

func MakeNewClientV1(apiKey string) (Client, error) {
	version := "v1"
	url, err := url.Parse(fmt.Sprintf("%s/%s/images/", baseUrl, version))
	if err != nil {
		return nil, fmt.Errorf("error parsing url: %s", err)
	}

	newClient := &ClientV1{
		version: version,
		client: &http.Client{
			Transport: &ApiKeyTransport{
				ApiKey:    apiKey,
				Transport: http.DefaultTransport,
			},
		},

		Url:       url,
		UserAgent: fmt.Sprintf("%s/%s", goDalle2, version),
	}

	return newClient, nil
}

func (c *ClientV1) makeRequest(ctx context.Context, path, contentType string, data interface{}) (req *http.Request, err error) {
	url, err := c.Url.Parse(path)
	if err != nil {
		return nil, err
	}

	switch contentType {
	case applicationJson:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(b)
		req, err = http.NewRequest("POST", url.String(), reader)
		if err != nil {
			return nil, err
		}
		req.Header.Set(headerContentType, applicationJson)
	case multipartForm:
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		if err := encodeFormData(writer, data); err != nil {
			return nil, err
		}
		req, err = http.NewRequest("POST", url.String(), bytes.NewReader(body.Bytes()))
		if err != nil {
			return nil, err
		}
		req.Header.Set(headerContentType, writer.FormDataContentType())
	default:
		return nil, fmt.Errorf("unsupported content-type: %s", contentType)
	}

	req.Header.Set(headerUserAgent, c.UserAgent)
	return req, err
}

func (c *ClientV1) handleResponse(resp *http.Response) (*Response, error) {
	if resp.StatusCode != 200 {
		errResp := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(errResp); err != nil {
			return nil, err
		}
		return nil, errResp
	}
	output := &Response{}
	if err := json.NewDecoder(resp.Body).Decode(output); err != nil {
		return nil, err
	}
	return output, nil
}

func (c *ClientV1) Create(ctx context.Context, prompt string, opts ...Option) (*Response, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	options := newDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	requestParams := &imageRequest{
		apiOption: options,
		Prompt:    prompt,
	}

	req, err := c.makeRequest(ctx, "generations", applicationJson, requestParams)
	if err != nil {
		return nil, err
	}
	req.Header[headerContentType] = []string{applicationJson}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.handleResponse(resp)
}

func (c *ClientV1) Edit(ctx context.Context, image, mask []byte, prompt string, opts ...Option) (*Response, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	options := newDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	requestParams := &editRequest{
		apiOption: options,
		Image:     image,
		Mask:      mask,
		Prompt:    prompt,
	}

	req, err := c.makeRequest(ctx, "edits", multipartForm, requestParams)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.handleResponse(resp)
}

func (c *ClientV1) Variation(ctx context.Context, image []byte, opts ...Option) (*Response, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	options := newDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	requestParams := &variationRequest{
		apiOption: options,
		Image:     image,
	}

	req, err := c.makeRequest(ctx, "variations", multipartForm, requestParams)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return c.handleResponse(resp)
}

func encodeFormData(writer *multipart.Writer, data interface{}) error {
	defer writer.Close()
	kv, err := query.Values(data)
	if err != nil {
		return err
	}
	strFields := []string{"n", "response_format", "size", "user", "format", "prompt"}
	byteFields := []string{"image", "mask"}

	// text fields
	for _, k := range strFields {
		if !kv.Has(k) {
			continue
		}
		if err := writer.WriteField(k, kv.Get(k)); err != nil {
			return err
		}
	}

	// image fields
	for _, k := range byteFields {
		if !kv.Has(k) {
			continue
		}
		fw, err := writer.CreateFormFile(k, k)
		if err != nil {
			return err
		}
		switch v := data.(type) {
		case *editRequest:
			switch k {
			case "image":
				if _, err = io.Copy(fw, bytes.NewReader(v.Image)); err != nil {
					return err
				}
			case "mask":
				if _, err = io.Copy(fw, bytes.NewReader(v.Mask)); err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid form field for request: %s", k)
			}
		case *variationRequest:
			switch k {
			case "image":
				if _, err = io.Copy(fw, bytes.NewReader(v.Image)); err != nil {
					return err
				}
			default:
				return fmt.Errorf("invalid form field for request: %s", k)
			}
		default:
			return errors.New("error casting request")
		}

	}
	return nil
}
