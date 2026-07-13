package sdkcore

import (
	"fmt"
	"net/http"
	"time"
)

// ClientConfig holds HTTP client settings for SDK consumers.
type ClientConfig struct {
	BaseURL string        `validate:"required,url"`
	APIKey  string        `validate:"required"`
	Timeout time.Duration `validate:"required"`
}

// HTTPClient is a foundation HTTP wrapper. Business methods added in Phase 2.
type HTTPClient struct {
	baseURL string
	apiKey  string
	http    *http.Client
}

// NewHTTPClient creates a configured SDK HTTP client.
func NewHTTPClient(cfg ClientConfig) (*HTTPClient, error) {
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &HTTPClient{
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		http:    &http.Client{Timeout: cfg.Timeout},
	}, nil
}

// Do executes an HTTP request with default SDK headers.
func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", "ai-finops-sdk-go/0.0.0")
	return c.http.Do(req)
}

// BaseURL returns the configured API base URL.
func (c *HTTPClient) BaseURL() string {
	return c.baseURL
}