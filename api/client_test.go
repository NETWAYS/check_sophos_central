package api_test

import (
	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/jarcoal/httpmock"
	"net/http"
	"os"
	"testing"
)

func envClient(t *testing.T) *api.Client {
	id := os.Getenv("SOPHOS_CLIENT_ID")
	secret := os.Getenv("SOPHOS_CLIENT_SECRET")

	if id == "" || secret == "" {
		t.Skip("SOPHOS_CLIENT_ID and SOPHOS_CLIENT_SECRET must be set!")
	}

	return api.NewClient(id, secret)
}

func testClient() (*api.Client, func()) {
	httpmock.Activate()

	httpmock.RegisterResponder("POST", "https://id.sophos.com/api/v2/oauth2/token",
		func(req *http.Request) (*http.Response, error) {
			body := map[string]interface{}{
				"access_token":  "<jwt>",
				"errorCode":     "success",
				"expires_in":    3600,
				"message":       "OK",
				"refresh_token": "<token>",
				"token_type":    "bearer",
				"trackingId":    "<uuid>",
			}
			return httpmock.NewJsonResponse(200, body)
		})

	httpmock.RegisterResponder("GET", "https://api.central.sophos.com/whoami/v1",
		func(req *http.Request) (*http.Response, error) {
			body := map[string]interface{}{
				"id":     "57ca9a6b-885f-4e36-95ec-290548c26059",
				"idType": "tenant",
				"apiHosts": map[string]string{
					"global":     "https://api.central.sophos.com",
					"dataRegion": "https://api-eu01.central.sophos.com",
				},
			}
			return httpmock.NewJsonResponse(200, body)
		})

	return api.NewClient("", ""), func() {
		httpmock.DeactivateAndReset()
	}
}
