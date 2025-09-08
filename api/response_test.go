package api_test

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"strings"
	"testing"
)

const testResponse = `{
    "pages": {
        "fromKey": "",
        "nextKey": "abcdef",
        "size": 50,
        "maxSize": 100
    },
    "items": [
		{"id": "c86dc437-daa8-4dec-b5a8-ce1e0a2c0c5e", "name": "test1"}
    ]
}`

const testResponseSecond = `{
    "pages": {
        "fromKey": "abcdef",
        "nextKey": "",
        "size": 50,
        "maxSize": 100
    },
    "items": [
		{"id": "6d497ca9-d5ef-495c-9074-a7021afc42c1", "name": "test2"}
    ]
}`

func TestClient_GetResults(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://api-eu01.central.sophos.com/common/v1/test",
		func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.RawQuery, "pageFromKey=abcdef") {
				return httpmock.NewStringResponse(200, testResponseSecond), nil
			}
			return httpmock.NewStringResponse(200, testResponse), nil
		})

	err := c.WhoAmI()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	req, err := c.NewDataRequest("GET", "common/v1/test", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = c.GetResults(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
