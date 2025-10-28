package fetcher

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ClientOption func(*FetcherClient)

type FetcherClient struct {
	client   *http.Client
	maxBytes int64
}

func NewClient(opts ...ClientOption) *FetcherClient {
	c := &FetcherClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxBytes: 10 * 1024 * 1024,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *FetcherClient) {
		c.client.Timeout = timeout
	}
}

func WithMaxBytes(maxBytes int64) ClientOption {
	return func(c *FetcherClient) {
		c.maxBytes = maxBytes
	}
}

func (c *FetcherClient) Fetch(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "NewsCatcher/1.0")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code %d", resp.StatusCode)
	}

	limitReader := io.LimitReader(resp.Body, c.maxBytes)
	return io.ReadAll(limitReader)
}
