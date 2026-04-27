package httpclient

import (
	"io"
	"net/http"
)

type Client struct {
	http *http.Client
}

func New() *Client {
	return &Client{http: &http.Client{}}
}

func (c *Client) Get(url string) ([]byte, error) {
	resp, err := c.http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
