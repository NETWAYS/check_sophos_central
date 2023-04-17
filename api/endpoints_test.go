package api_test

import (
	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

const endpointsResponse = `{
    "pages": {
        "fromKey": "",
        "nextKey": "",
        "size": 50,
        "maxSize": 100
    },
    "items": [
        {
            "id": "4f725e46-a174-45a7-a352-49ed15840ab1",
            "type": "computer",
            "tenant": {
                "id": "2a591e9c-5a0b-4623-9e78-f28692a1616f"
            },
            "hostname": "DESKTOP-F82003",
            "health": {
                "overall": "good",
                "threats": {
                    "status": "good"
                },
                "services": {
                    "status": "bad",
                    "serviceDetails": [
                        {
                            "name": "HitmanPro.Alert service",
                            "status": "stopped"
                        },
                        {
                            "name": "Sophos Anti-Virus",
                            "status": "running"
                        },
                        {
                            "name": "Sophos Anti-Virus Status Reporter",
                            "status": "running"
                        }
                    ]
                }
            },
            "os": {
                "isServer": false,
                "platform": "windows",
                "name": "Windows 10 Pro",
                "majorVersion": 10,
                "minorVersion": 0,
                "build": 17134
            },
            "ipv4Addresses": [
                "10.10.10.100"
            ],
            "ipv6Addresses": [
                "fe80::9f3e:aeff:8892:0bb8"
            ],
            "macAddresses": [
                "00:0B:FF:E1:C2:C2"
            ],
            "associatedPerson": {
                "viaLogin": "DESKTOP-F82003\\admin"
            },
            "assignedProducts": [
                {
                    "code": "endpointProtection",
                    "version": "10.8.3.441"
                },
                {
                    "code": "deviceEncryption",
                    "version": "1.4.103"
                },
                {
                    "code": "interceptX",
                    "version": "2.0.14"
                },
                {
                    "code": "coreAgent",
                    "version": "2.4.0"
                }
            ],
            "lastSeenAt": "2019-08-29T10:20:29.375Z"
        }
    ]
}`

func TestEndpoint_String(t *testing.T) {
	testcases := map[string]struct {
		endpoint api.Endpoint
		expected string
	}{
		"default": {
			endpoint: api.Endpoint{},
			expected: " []  ",
		},
		"custom-endpoint": {
			endpoint: api.Endpoint{
				Hostname: "hostname",
				ID:       "ID",
				Health: api.EndpointHealth{
					Overall: "cataclysmic",
				},
				Type: "Typ",
			},
			expected: "hostname [cataclysmic] ID Typ",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.endpoint.String())
		})
	}
}

func TestClient_GetEndpoints(t *testing.T) {
	c := envClient(t)

	err := c.WhoAmI()
	assert.NoError(t, err)

	_, err = c.GetEndpoints()
	assert.NoError(t, err)
}

func TestClient_GetEndpoints_Mock(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	httpmock.RegisterResponder("GET", "https://api-eu01.central.sophos.com/endpoint/v1/endpoints",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, endpointsResponse), nil
		})

	err := c.WhoAmI()
	assert.NoError(t, err)

	_, err = c.GetEndpoints()
	assert.NoError(t, err)
}
