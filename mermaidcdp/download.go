package mermaidcdp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DownloadJSSource downloads and returns the specified version of MermaidJS from a CDN.
// Use this if you cannot bundle MermaidJS with your application.
//
// version is the version of MermaidJS to download, without a "v" prefix.
// Example values are "10", "10.3.0".
func DownloadJSSource(ctx context.Context, version string) (string, error) {
	var d downloader
	return d.Download(ctx, version)
}

type downloader struct {
	httpClient *http.Client // for testing
}

func (d *downloader) Download(ctx context.Context, version string) (string, error) {
	if version == "" {
		return "", fmt.Errorf("verison must be specified")
	}

	httpClient := http.DefaultClient
	if d.httpClient != nil {
		httpClient = d.httpClient
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://cdn.jsdelivr.net/npm/mermaid@"+version+"/dist/mermaid.min.js",
		nil, // body
	)
	if err != nil {
		return "", fmt.Errorf("prepare request: %w", err)
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, res.Body) // drain response body
		return "", fmt.Errorf("http status %d", res.StatusCode)
	}

	var s strings.Builder
	if _, err := io.Copy(&s, res.Body); err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}
	return s.String(), nil
}
