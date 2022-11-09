package dalle2

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

type mockRoundTripper struct {
	TestKey string
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if v := req.Header.Get(headerAuth); !strings.Contains(v, m.TestKey) {
		return nil, errors.New("invalid auth header")
	}
	return &http.Response{}, nil
}

func TestRoundTrip(t *testing.T) {
	testTransport := &ApiKeyTransport{
		ApiKey: "test",
		Transport: &mockRoundTripper{
			TestKey: "test",
		},
	}

	if _, err := testTransport.RoundTrip(&http.Request{
		Header: map[string][]string{
			headerAuth: {"test"},
		},
	}); err != nil {
		t.Error(err)
	}
}
