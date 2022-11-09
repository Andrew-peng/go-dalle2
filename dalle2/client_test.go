package dalle2

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testApiKey = "thisisatestkey"
const testVersion = "vtest"

func setup() (client *ClientV1, cleanup func()) {
	cl, _ := MakeNewClientV1(testApiKey)
	client = cl.(*ClientV1)
	apiHandler := http.NewServeMux()
	// Create
	apiHandler.HandleFunc(fmt.Sprintf("/%s/images/generations", testVersion), func(w http.ResponseWriter, r *http.Request) {
		req := &imageRequest{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set(headerContentType, applicationJson)
		if req.User == "return_err" {
			resp := &ErrorResponse{
				ErrorDetails: &ErrorDetails{
					Code:    404,
					Message: "this is an error",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := &Response{
				Data: []*OutputImage{
					{
						Url: "testurl",
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}
	})

	// Edit
	apiHandler.HandleFunc(fmt.Sprintf("/%s/images/edits", testVersion), func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1064); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		formData := r.MultipartForm
		w.Header().Set(headerContentType, applicationJson)
		if user, ok := formData.Value["user"]; ok && user[0] == "return_err" {
			resp := &ErrorResponse{
				ErrorDetails: &ErrorDetails{
					Code:    404,
					Message: "this is an error",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := &Response{
				Data: []*OutputImage{
					{
						Url: "testurl",
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}
	})

	// Variations
	apiHandler.HandleFunc(fmt.Sprintf("/%s/images/variations", testVersion), func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(1064); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		formData := r.MultipartForm
		w.Header().Set(headerContentType, applicationJson)
		if user, ok := formData.Value["user"]; ok && user[0] == "return_err" {
			resp := &ErrorResponse{
				ErrorDetails: &ErrorDetails{
					Code:    404,
					Message: "this is an error",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
		} else {
			resp := &Response{
				Data: []*OutputImage{
					{
						Url: "testurl",
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		}
	})
	server := httptest.NewServer(apiHandler)
	serverUrl, _ := url.Parse(fmt.Sprintf("%s/%s/images/", server.URL, testVersion))
	client.Url = serverUrl
	cleanup = func() {
		server.Close()
	}
	return client, cleanup
}

func TestCreate(t *testing.T) {
	client, cleanup := setup()
	defer cleanup()
	resp, err := client.Create(context.Background(), "this is a test prompt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, resp.Data, 1)
	assert.EqualValues(t, "testurl", resp.Data[0].Url)

	resp, err = client.Create(context.Background(), "this is a test prompt", WithUser("return_err"))
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestEdit(t *testing.T) {
	client, cleanup := setup()
	defer cleanup()
	resp, err := client.Edit(context.Background(), []byte("image"), []byte("mask"), "this is a test prompt")
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, resp.Data, 1)
	assert.EqualValues(t, "testurl", resp.Data[0].Url)

	resp, err = client.Edit(context.Background(), []byte("image"), []byte("mask"), "this is a test prompt", WithUser("return_err"))
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}

func TestVariation(t *testing.T) {
	client, cleanup := setup()
	defer cleanup()
	resp, err := client.Variation(context.Background(), []byte("image"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, resp.Data, 1)
	assert.EqualValues(t, "testurl", resp.Data[0].Url)

	resp, err = client.Variation(context.Background(), []byte("image"), WithUser("return_err"))
	assert.Nil(t, resp)
	assert.NotNil(t, err)
}
