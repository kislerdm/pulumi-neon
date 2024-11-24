package telemetry

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPClient_setUAHeader(t *testing.T) {
	tests := []struct {
		name   string
		c      *HTTPClient
		wantUA string
	}{
		{
			name:   "provider Foo@1.0.0",
			c:      NewHTTPClient("Foo", "1.0.0"),
			wantUA: "pulumiProvider-Foo@1.0.0",
		},
		{
			name:   "provider version is unknown",
			c:      NewHTTPClient("Bar", ""),
			wantUA: "",
		},
		{
			name:   "provider name is unknown",
			c:      NewHTTPClient("Bar", ""),
			wantUA: "",
		},
	}

	t.Parallel()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.c)
			var r = &http.Request{
				URL: &url.URL{
					Scheme: "https",
					Host:   "foo.com",
				},
			}
			tt.c.setUAHeader(r)
			assert.Equal(t, tt.wantUA, r.Header.Get("User-Agent"))
		})
	}
}
