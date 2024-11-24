// Package telemetry defines the HTTP client to ensure the API calls are instrumented to include required headers.
package telemetry

import (
	"fmt"
	"net/http"
	"time"
)

// NewHTTPClient init the HTTP client to send HTTP request with required telemetry headers.
func NewHTTPClient(providerName, providerVersion string) *HTTPClient {
	return &HTTPClient{
		ProviderName:    providerName,
		ProviderVersion: providerVersion,
		c:               &http.Client{Timeout: 2 * time.Minute},
	}
}

// HTTPClient instrumented HTTP client.
type HTTPClient struct {
	ProviderName    string
	ProviderVersion string

	c *http.Client
}

func (c HTTPClient) Do(r *http.Request) (*http.Response, error) {
	c.setUAHeader(r)

	return c.c.Do(r)
}

func (c HTTPClient) setUAHeader(r *http.Request) {
	if c.ProviderName != "" && c.ProviderVersion != "" {
		if r.Header == nil {
			r.Header = make(http.Header)
		}

		telemetryHeader := fmt.Sprintf("pulumiProvider-%s@%s", c.ProviderName, c.ProviderVersion)

		r.Header.Set("User-Agent", telemetryHeader)
	}
}
