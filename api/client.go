package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	DefaultURL      = "https://api.central.sophos.com"
	AuthTokenURL    = "https://id.sophos.com/api/v2/oauth2/token" //nolint:gosec
	DefaultTimeout  = 5
	DefaultPageSize = 100
)

type Client struct {
	AuthConfig *clientcredentials.Config
	HTTPClient *http.Client
	BaseURL    string
	DataURL    string
	UserInfo   *UserInfo
	TenantID   string
	PageSize   uint32
}

func NewClient(id, secret string) (c *Client) {
	c = &Client{
		BaseURL:  DefaultURL,
		PageSize: DefaultPageSize,
	}

	// Prepare custom client that using a logging transport.
	client := http.DefaultClient
	client.Transport = LoggingRoundTripper{http.DefaultTransport}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)

	// Setup authentication.
	c.AuthConfig = &clientcredentials.Config{
		ClientID:       id,
		ClientSecret:   secret,
		TokenURL:       AuthTokenURL,
		AuthStyle:      oauth2.AuthStyleInParams,
		EndpointParams: url.Values{"scope": {"token"}},
	}
	c.HTTPClient = c.AuthConfig.Client(ctx)
	c.HTTPClient.Timeout = time.Duration(DefaultTimeout) * time.Second

	return
}

func (c *Client) NewRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(context.Background(), method, c.BaseURL+"/"+url, body)
	if err != nil {
		err = fmt.Errorf("could not create http request: %w", err)
	}

	return
}

func (c *Client) NewDataRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	if c.DataURL == "" {
		err = errors.New("DataURL is not set, call whoami first")
		return
	}

	if c.TenantID == "" {
		err = errors.New("TenantID is not configured")
		return
	}

	req, err = http.NewRequestWithContext(context.Background(), method, c.DataURL+"/"+url, body)

	if err != nil {
		err = fmt.Errorf("could not create http request: %w", err)
		return
	}

	req.Header.Set("X-Tenant-Id", c.TenantID)

	return
}

func (c *Client) Do(req *http.Request) (res *http.Response, err error) {
	res, err = c.HTTPClient.Do(req)
	if err != nil {
		err = fmt.Errorf("HTTP request failed: %w", err)
	}

	return
}
