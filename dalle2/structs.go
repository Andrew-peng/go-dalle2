package dalle2

import (
	"fmt"
)

const (
	SMALL  = "256x256"
	MEDIUM = "512x512"
	LARGE  = "1024x1024"

	URL    = "url"
	BASE64 = "b64_json"
)

type apiOption struct {
	N      uint8  `json:"n,omitempty" url:"n,omitempty"`
	Size   string `json:"size,omitempty" url:"size,omitempty"`
	Format string `json:"response_format,omitempty" url:"response_format,omitempty"`
	User   string `json:"user,omitempty" url:"user,omitempty"`
}

type imageRequest struct {
	apiOption
	Prompt string `json:"prompt,omitempty"`
}

type editRequest struct {
	apiOption
	Image  []byte `url:"image,omitempty"`
	Mask   []byte `url:"mask,omitempty"`
	Prompt string `url:"prompt,omitempty"`
}

type variationRequest struct {
	apiOption
	Image []byte `url:"image,omitempty"`
}

type Option func(*apiOption)

func newDefaultOptions() apiOption {
	return apiOption{
		N:      1,
		Size:   LARGE,
		Format: URL,
		User:   "",
	}
}

func WithNumImages(n uint8) Option {
	return func(o *apiOption) {
		o.N = n
	}
}

func WithSize(size string) Option {
	return func(o *apiOption) {
		o.Size = size
	}
}

func WithFormat(format string) Option {
	return func(o *apiOption) {
		o.Format = format
	}
}

func WithUser(user string) Option {
	return func(o *apiOption) {
		o.User = user
	}
}

type OutputImage struct {
	Base64 string `json:"b64_json,omitempty"`
	Url    string `json:"url,omitempty"`
}

type Response struct {
	Created int64          `json:"created,omitempty"`
	Data    []*OutputImage `json:"data,omitempty"`
}

type ErrorDetails struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Param   string `json:"param,omitempty"`
	Type    string `json:"type,omitempty"`
}

type ErrorResponse struct {
	ErrorDetails *ErrorDetails `json:"error,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s: %s", e.ErrorDetails.Type, e.ErrorDetails.Message)
}
