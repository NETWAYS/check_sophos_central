package api

import (
	"encoding/json"
	"fmt"
)

type Endpoint struct {
	// The hostname of the endpoint.
	Hostname string `json:"hostname"`
	// The unique ID for the endpoint.
	ID string `json:"id"`
	// Products assigned to the endpoint.
	AssignedProducts []json.RawMessage `json:"assignedProducts"`
	// The health status of the endpoint.
	Health EndpointHealth `json:"health"`
	// The endpoint type.
	//
	// The following values are allowed:
	// computer, server, securityVm
	Type string `json:"type"`
}

type EndpointHealth struct {
	// Health status of the endpoint or a service running on the endpoint.
	//
	// The following values are allowed:
	// good, suspicious, bad, unknown
	Overall string `json:"overall"`
	// Status of services on the endpoint.
	Services json.RawMessage `json:"services"`
	// Threats on the endpoint.
	Threats json.RawMessage `json:"threats"`
}

func (e *Endpoint) String() string {
	return fmt.Sprintf(
		"%v [%v] %v %v",
		e.Hostname,
		e.Health.Overall,
		e.ID,
		e.Type,
	)
}

func (c *Client) GetEndpoints() (endpoints []*Endpoint, err error) {
	req, err := c.NewDataRequest("GET", "endpoint/v1/endpoints", nil)
	if err != nil {
		return
	}

	items, err := c.GetResults(req)
	if err != nil {
		return
	}

	// retrieve items from response
	for _, item := range items {
		e := &Endpoint{}

		err = json.Unmarshal(item, e)
		if err != nil {
			return
		}

		endpoints = append(endpoints, e)
	}

	return
}
