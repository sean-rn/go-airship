package airship

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sean-rn/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestBearerToken = "test-ua-token"

// This suite tests the Urban Airship HTTP client itself,
// verifying that it sends the correct JSON request bodies to the right endpoints.

func TestInvokeEndpoint_BasicAuth(t *testing.T) {
	assert := assert.New(t)

	expectedBody := `{ "message": "Hello World" }`

	client := httpmock.NewHandlerClient(func(rw http.ResponseWriter, req *http.Request) {
		// Validate request parameters and body
		assert.Equal("POST", req.Method)
		assert.Equal("https://go.urbanairship.com/api/push", req.URL.String())
		assert.Equal("application/json", req.Header.Get("Content-Type"))
		assert.Equal("Basic YXBwLWtleTptYXN0ZXItc2VjcmV0", req.Header.Get("Authorization"))
		assert.Equal("application/vnd.urbanairship+json; version=3;", req.Header.Get("Accept"))
		assertBodyJSONEqual(t, expectedBody, req.Body)
		// Write response
		rw.Write([]byte(`{"ok": true,"operation_id": "df6a6b50","push_ids": ["9d78a53b"],"message_ids": [], "content_urls": []}`))
	})

	testConnection := New(WithHTTPClient(client), WithBasicAuth("app-key", "master-secret"))

	// Invoke!
	body := map[string]string{"message": "Hello World"}
	err := testConnection.InvokeEndpoint(http.MethodPost, "/api/push", body)
	require.Nil(t, err)
}

func TestInvokeEndpoint_BearerToken(t *testing.T) {
	assert := assert.New(t)

	expectedBody := `{ "message": "Hello World" }`

	client := httpmock.NewHandlerClient(func(rw http.ResponseWriter, req *http.Request) {
		// Validate request parameters and body
		assert.Equal("PUT", req.Method)
		assert.Equal("https://go.urbanairship.com/api/other", req.URL.String())
		assert.Equal("application/json", req.Header.Get("Content-Type"))
		assert.Equal("Bearer test-ua-token", req.Header.Get("Authorization"))
		assert.Equal("application/vnd.urbanairship+json; version=3;", req.Header.Get("Accept"))
		assertBodyJSONEqual(t, expectedBody, req.Body)
		// Write response
		rw.Write([]byte(`{"ok": true,"operation_id": "df6a6b50","push_ids": ["9d78a53b"],"message_ids": [], "content_urls": []}`))
	})

	testConnection := New(WithHTTPClient(client), WithBearerAuth(TestBearerToken))

	// Invoke!
	body := map[string]string{"message": "Hello World"}
	err := testConnection.InvokeEndpoint(http.MethodPut, "/api/other", body)
	require.Nil(t, err)
}

// Make sure HTTP error codes returned by httpClient.Do() don't panic
func TestInvokeEndpoint_HttpError(t *testing.T) {
	assert := assert.New(t)

	client := httpmock.NewHandlerClient(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusForbidden)
		rw.Write([]byte(`{"error": "Forbidden"}`))
	})

	testConnection := New(WithBearerAuth(TestBearerToken), WithHTTPClient(client))

	// Invoke!
	body := map[string]string{"message": "Hello World"}
	err := testConnection.InvokeEndpoint(http.MethodPost, "/api/push", body)
	assert.Error(err) // HTTP errors return a go error
}

// Make sure errors returned by httpClient.Do() don't panic and are returned.
func TestInvokeEndpoint_GoError(t *testing.T) {
	assert := assert.New(t)

	client := httpmock.NewTransportClient(func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("oh no an error")
	})

	testConnection := New(WithBearerAuth(TestBearerToken), WithHTTPClient(client))

	// Invoke!
	body := map[string]string{"message": "Hello World"}
	err := testConnection.InvokeEndpoint(http.MethodPost, "/api/push", body)
	assert.Error(err)
}

// Helper to read and compare the request body.
func assertBodyJSONEqual(t testing.TB, expected string, body io.ReadCloser, msgAndArgs ...interface{}) bool {
	// Read the body and check for error while reading.
	defer body.Close()
	bodyBytes, err := ioutil.ReadAll(body)
	require.Nil(t, err, msgAndArgs...)
	return assert.JSONEq(t, expected, string(bodyBytes), msgAndArgs...)
}
