package api_test

import (
	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/jarcoal/httpmock"
	"net/http"
	"testing"
	"time"
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

func TestAlert_String(t *testing.T) {
	testcases := map[string]struct {
		alert    api.Alert
		expected string
	}{
		"default": {
			alert: api.Alert{
				RaisedAt: time.Date(2009, time.November, 11, 23, 0, 0, 0, time.UTC),
			},
			expected: "2009-11-11 23:00 []      ",
		},
		"custom-alert": {
			alert: api.Alert{
				ID:          "ID1",
				RaisedAt:    time.Date(2009, time.November, 11, 23, 0, 0, 0, time.UTC),
				Category:    "Cat",
				Description: "Desc",
				GroupKey:    "GK",
				Product:     "Product",
				Severity:    "critical",
				Type:        "type",
			},
			expected: "2009-11-11 23:00 [critical] ID1 Product Desc type Cat GK",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			if tc.expected != tc.alert.String() {
				t.Fatalf("expected %v, got %v", tc.expected, tc.alert.String())
			}
		})
	}
}

func TestClient_GetAlerts(t *testing.T) {
	c := envClient(t)

	err := c.WhoAmI()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = c.GetAlerts()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestClient_GetAlerts_Mock(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://api-eu01.central.sophos.com/common/v1/alerts",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, alertsResponse), nil
		})

	err := c.WhoAmI()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = c.GetAlerts()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
