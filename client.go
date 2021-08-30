package airship

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	// BaseURL is the base of the API endpoints https://docs.airship.com/api/ua/#servers
	BaseURL = "https://go.urbanairship.com"
	// AcceptHeader is the value to send in the Accept header as required by docs.
	AcceptHeader = "application/vnd.urbanairship+json; version=3;"
)

const (
	// EndpointPushToTemplate is the path of the "/api/templates/push" POST endpoint.
	// https://docs.airship.com/api/ua/#operation-api-templates-push-post
	EndpointPushToTemplate = "/api/templates/push"
	// EndpointSendPush is the path of the "Send a Push" POST endpoint.
	// https://docs.airship.com/api/ua/#operation-api-push-post
	EndpointSendPush = "/api/push"
	// EndpointCreateAndSend is the path of the "Create and Send" POST endpoint.
	// https://docs.airship.com/api/ua/#operation-api-create-and-send-post
	EndpointCreateAndSend = "/api/create-and-send"
)

//go:generate mockery --name Client
// Install mockery from https://github.com/vektra/mockery

// Client is the API for interacting with Urban Airship
type Client interface {
	InvokeEndpoint(method string, endpoint string, body interface{}) error
}

// Urban Airship HTTP API Client implementation
type uaHTTPClient struct {
	httpClient  *http.Client
	authHeader  string
	endpointURL string
}

// ClientOption are configuration functions that can be passed to New to configure the client.
type ClientOption func(c *uaHTTPClient)

// New creates a new client instance configured with the given options.
// For example:
//    conn := airship.New(WithBasicAuth("app-key", "master-secret"))
func New(options ...ClientOption) Client {
	client := uaHTTPClient{}
	for _, opt := range options {
		opt(&client)
	}
	if client.httpClient == nil {
		client.httpClient = &http.Client{}
	}
	if client.endpointURL == "" {
		client.endpointURL = BaseURL
	}
	return &client
}

// WithBasicAuth configures the Airship Client to use HTTP Basic Auth
// https://docs.airship.com/api/ua/#security-basicauth
func WithBasicAuth(appKey, masterSecret string) ClientOption {
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(appKey+":"+masterSecret))
	return func(c *uaHTTPClient) {
		c.authHeader = authHeader
	}
}

// WithBearerAuth configures the Airship Client to use Bearer Auth
// https://docs.airship.com/api/ua/#security-bearerauth
func WithBearerAuth(token string) ClientOption {
	authHeader := "Bearer " + token
	return func(c *uaHTTPClient) {
		c.authHeader = authHeader
	}
}

// WithHTTPClient overrides the http.Client instance used by the Airship Client.
// This is useful for unit tests of the client itself, but not much else.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *uaHTTPClient) {
		c.httpClient = httpClient
	}
}

// InvokeEndpoint invokes the airship API endpoint by sending <body> to <endpoint> using HTTP <method>.
// The response body is discarded unless an error status is returned.
func (cfg *uaHTTPClient) InvokeEndpoint(method string, endpoint string, body interface{}) error {
	jsonStr, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, cfg.endpointURL+endpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", cfg.authHeader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", AcceptHeader)

	resp, err := cfg.httpClient.Do(req)
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			respBody, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("airship: request returned %d: %s", resp.StatusCode, respBody)
		}
	}
	return err
}
