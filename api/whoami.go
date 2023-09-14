package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type UserInfo struct {
	ID       string            `json:"id"`
	IDType   string            `json:"idType"`
	APIHosts map[string]string `json:"apiHosts"`
}

func (c *Client) WhoAmI() (err error) {
	req, err := c.NewRequest(http.MethodGet, "whoami/v1", nil) //nolint: noctx
	if err != nil {
		return
	}

	resp, err := c.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("whoami request failed with status: %s", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	info := &UserInfo{}

	err = json.Unmarshal(body, info)
	if err != nil {
		return
	}

	c.UserInfo = info

	// parse and set additional API endpoints
	if val, ok := info.APIHosts["global"]; ok && c.BaseURL != val {
		c.BaseURL = val
	}

	if val, ok := info.APIHosts["dataRegion"]; ok {
		c.DataURL = val
	} else {
		err = fmt.Errorf("missing dataRegion value under apiHosts in whoami: %s", string(body))
		return
	}

	// set TenantID when Token belongs to a tenant
	if info.IDType == "tenant" && c.TenantID == "" {
		c.TenantID = info.ID
	}

	return
}
