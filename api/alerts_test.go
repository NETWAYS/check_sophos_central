package api_test

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const alertsResponse = `{
    "pages": {
        "fromKey": "",
        "nextKey": "",
        "size": 50,
        "maxSize": 100
    },
    "items": [
        {
			"id": "eefa2158-ca7b-4150-8a3b-1fe06cd808f5",
			"allowedActions": [],
			"category": "systemHealth",
            "description": "Some fancy alerts",
            "groupKey": "group",
            "managedAgent": {
				"id": "d36e723b-9c04-421c-96fb-4c691ebce985",
				"type": "computer"
			},
            "person": {"id": "41471d96-fbe0-479c-ba13-ae09e575a1fc"},
            "product": "endpoint",
            "raisedAt": "2019-08-29T10:20:29.375Z",
            "severity": "medium",
            "tenant": {
				"id": "57ca9a6b-885f-4e36-95ec-290548c26059",
				"name": "Mustermann"
             },
            "type": "type"
		}
	]
}`

func TestClient_GetAlerts(t *testing.T) {
	c := envClient(t)

	err := c.WhoAmI()
	assert.NoError(t, err)

	_, err = c.GetAlerts()
	assert.NoError(t, err)
}

func TestClient_GetAlerts_Mock(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://api-eu01.central.sophos.com/common/v1/alerts",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, alertsResponse), nil
		})

	err := c.WhoAmI()
	assert.NoError(t, err)

	_, err = c.GetAlerts()
	assert.NoError(t, err)
}
