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
	"os"

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

var (
	errNilCtx = errors.New("received nil context")
)

type Client interface {
	Create(context.Context, string, ...Option) (*Response, error)
	Edit(context.Context, string, string, string, ...Option) (*Response, error)
	Variation(context.Context, string, ...Option) (*Response, error)
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
	if ctx == nil {
		return nil, errNilCtx
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

func (c *ClientV1) Edit(ctx context.Context, imagePath, maskPath, prompt string, opts ...Option) (*Response, error) {
	if ctx == nil {
		return nil, errNilCtx
	}

	options := newDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	requestParams := &editRequest{
		apiOption: options,
		Image:     imagePath,
		Mask:      maskPath,
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

func (c *ClientV1) Variation(ctx context.Context, imagePath string, opts ...Option) (*Response, error) {
	if ctx == nil {
		return nil, errNilCtx
	}

	options := newDefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	requestParams := &variationRequest{
		apiOption: options,
		Image:     imagePath,
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

	// text fields
	for _, k := range []string{"n", "response_format", "size", "user", "format", "prompt"} {
		if !kv.Has(k) {
			continue
		}
		if err := writer.WriteField(k, kv.Get(k)); err != nil {
			return err
		}
	}

	// image fields
	for _, k := range []string{"image", "mask"} {
		if !kv.Has(k) {
			continue
		}
		fw, err := writer.CreateFormFile(k, kv.Get(k))
		if err != nil {
			return err
		}
		img, err := os.Open(kv.Get(k))
		if err != nil {
			return err
		}
		defer img.Close()
		if _, err = io.Copy(fw, img); err != nil {
			return err
		}
	}
	return nil
}
