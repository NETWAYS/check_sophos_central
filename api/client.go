package api

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultURL     = "https://api.central.sophos.com"
	AuthTokenURL   = "https://id.sophos.com/api/v2/oauth2/token" //nolint:gosec
	DefaultTimeout = 5
)

type Client struct {
	AuthConfig *clientcredentials.Config
	HttpClient *http.Client
	BaseURL    string
	DataURL    string
	UserInfo   *UserInfo
	TenantID   string
}

func NewClient(id, secret string) (c *Client) {
	c = &Client{
		BaseURL: DefaultURL,
	}

	ctx := context.Background()
	c.AuthConfig = &clientcredentials.Config{
		ClientID:       id,
		ClientSecret:   secret,
		TokenURL:       AuthTokenURL,
		AuthStyle:      oauth2.AuthStyleInParams,
		EndpointParams: url.Values{"scope": {"token"}},
	}
	c.HttpClient = c.AuthConfig.Client(ctx)
	c.HttpClient.Timeout = time.Duration(DefaultTimeout) * time.Second

	return
}

func (c *Client) GetCommonURL(urlPart string) string {
	if urlPart[0] != '/' {
		urlPart = "/" + urlPart
	}

	return c.BaseURL + urlPart
}

func (c *Client) NewDataRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	if c.DataURL == "" {
		err = fmt.Errorf("DataURL is not set, call whoami first")
		return
	}
	if c.TenantID == "" {
		err = fmt.Errorf("TenantID is not configured")
		return
	}

	req, err = http.NewRequest("GET", c.DataURL+"/"+url, body)
	if err != nil {
		return
	}

	req.Header.Set("X-Tenant-ID", c.TenantID)

	return
}
