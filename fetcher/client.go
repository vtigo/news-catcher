package fetcher

import (
	"io"
	"net/http"
	"time"
)

type Fetcher interface {
	Fetch(url string) ([]byte, error)
}

type ClientOption func(*FetcherClient)

type FetcherClient struct {
	client *http.Client
}

func NewClient(opts ...ClientOption) *FetcherClient {
	c := &FetcherClient{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
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

func (c *FetcherClient) Fetch(url string) ([]byte, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
