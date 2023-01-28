package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var (
	timeout = 10 * time.Second
)

func fetch(ctx context.Context, source *url.URL) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, source.String(), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("HTTP status error: %s", resp.Status)
	}

	return resp.Body, nil
}
