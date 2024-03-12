package blockchain

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	url        *url.URL
	httpClient *http.Client
}

func New(u *url.URL, httpClient *http.Client) *Client {
	return &Client{
		url:        u,
		httpClient: httpClient,
	}
}

const (
	createPoolPath        = "tokens/pools"
	mintTokenPath         = "tokens/mint"
	transferTokenPath     = "tokens/transfers"
	applicationJsonHeader = "application/json"
)

func (c *Client) CreatePool(ctx context.Context) error {
	p := c.url.JoinPath(createPoolPath)
	body := []byte(`{
		"name": "kaleido",
		"type": "nonfungible"
	  }`)
	_, err := c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", p, err)
	}
	return nil
}

func (c *Client) MintToken(ctx context.Context) error {
	p := c.url.JoinPath(mintTokenPath)
	body := []byte(`{
		"pool": "kaleido",
		"amount": "1"
	  }`)
	_, err := c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", p, err)
	}
	// TODO: store IDs into DB
	return nil
}

func (c *Client) TransferToken(ctx context.Context) error {
	p := c.url.JoinPath(transferTokenPath)
	// TODO: recipient address should be dynamic
	body := []byte(`{
		pool: 'kaleido',
		to: '0x2c499018bb8bc8a56065c07daeff4bceec254928',
		tokenIndex: '1',
		amount: '1',
	  }`)
	_, err := c.httpClient.Post(p.String(), applicationJsonHeader, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("c.httpClient.Post to (%s): %w", p, err)
	}
	// TODO: store IDs into DB
	return nil
}
