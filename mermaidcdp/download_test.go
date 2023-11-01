package mermaidcdp

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_DownloadMermaidJS(t *testing.T) {
	t.Parallel()

	giveBody := `const mermaid = {
		render: function() {
			return "<svg>...</svg>";
		},
	}`

	roundTrip := func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "GET", req.Method)
		assert.Equal(t,
			"https://cdn.jsdelivr.net/npm/mermaid@10/dist/mermaid.min.js",
			req.URL.String(),
		)

		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(giveBody)),
		}, nil
	}

	d := downloader{
		httpClient: &http.Client{
			Transport: &httpRoundTripper{RoundTripF: roundTrip},
		},
	}

	ctx := context.Background()
	got, err := d.Download(ctx, "10")
	require.NoError(t, err)
	assert.Equal(t, giveBody, got)
}

func TestConfig_DownloadMermaidJS_badVersions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		version string
		wantErr string
	}{
		{
			name:    "empty",
			wantErr: "verison must be specified",
		},
		{
			name:    "malformed URL",
			version: "foo\nbar",
			wantErr: "invalid control character in URL",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			roundTrip := func(req *http.Request) (*http.Response, error) {
				t.Fatalf("Unexpected HTTP request: %v %v", req.Method, req.URL)
				return nil, errors.New("unexpected HTTP request")
			}
			d := downloader{
				httpClient: &http.Client{
					Transport: &httpRoundTripper{RoundTripF: roundTrip},
				},
			}

			_, err := d.Download(context.Background(), tt.version)
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestConfig_DownloadMermaidJS_httpErrors(t *testing.T) {
	t.Parallel()

	t.Run("request failure", func(t *testing.T) {
		t.Parallel()

		roundTrip := func(*http.Request) (*http.Response, error) {
			return nil, errors.New("great sadness")
		}
		d := downloader{
			httpClient: &http.Client{
				Transport: &httpRoundTripper{RoundTripF: roundTrip},
			},
		}

		_, err := d.Download(context.Background(), "10")
		assert.ErrorContains(t, err, "great sadness")
	})

	t.Run("non 200 response", func(t *testing.T) {
		t.Parallel()

		roundTrip := func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader("404 not found")),
			}, nil
		}
		d := downloader{
			httpClient: &http.Client{
				Transport: &httpRoundTripper{RoundTripF: roundTrip},
			},
		}

		_, err := d.Download(context.Background(), "10")
		assert.ErrorContains(t, err, "http status 404")
	})

	t.Run("response body read error", func(t *testing.T) {
		t.Parallel()

		roundTrip := func(*http.Request) (*http.Response, error) {
			r := iotest.ErrReader(errors.New("great sadness"))
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(r),
			}, nil
		}
		d := downloader{
			httpClient: &http.Client{
				Transport: &httpRoundTripper{RoundTripF: roundTrip},
			},
		}

		_, err := d.Download(context.Background(), "10")
		assert.ErrorContains(t, err, "great sadness")
	})
}

type httpRoundTripper struct {
	RoundTripF func(req *http.Request) (*http.Response, error)
}

var _ http.RoundTripper = (*httpRoundTripper)(nil)

func (rt *httpRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.RoundTripF(req)
}
