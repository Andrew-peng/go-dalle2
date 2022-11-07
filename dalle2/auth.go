package dalle2

import (
	"errors"
	"fmt"
	"net/http"
)

type ApiKeyTransport struct {
	ApiKey    string
	Transport http.RoundTripper
}

func (t *ApiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ApiKey == "" {
		return nil, errors.New("empty API Key")
	}
	copyReq := &http.Request{}
	*copyReq = *req
	copyReq.Header = make(http.Header)
	for k, v := range req.Header {
		copyReq.Header[k] = v
	}
	copyReq.Header[headerAuth] = []string{fmt.Sprintf("Bearer %s", t.ApiKey)}
	return t.Transport.RoundTrip(copyReq)
}
