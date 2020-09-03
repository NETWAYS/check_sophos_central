package main

import (
	"fmt"
	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/NETWAYS/go-check"
)

type EndpointOverview struct {
	Total      int
	Good       []string
	Suspicious []string
	Bad        []string
	Unknown    []string
}

func CheckEndpoints(client *api.Client) (o *EndpointOverview, err error) {
	o = &EndpointOverview{}

	endpoints, err := client.GetEndpoints()
	if err != nil {
		return
	}

	for _, endpoint := range endpoints {
		o.Total++

		switch endpoint.Health.Overall {
		case "good":
			o.Good = append(o.Good, endpoint.Hostname)
		case "suspicious":
			o.Suspicious = append(o.Suspicious, endpoint.Hostname)
		case "bad":
			o.Bad = append(o.Bad, endpoint.Hostname)
		default:
			o.Unknown = append(o.Unknown, endpoint.Hostname)
		}
	}

	return
}
func (o *EndpointOverview) GetSummary() (s string) {
	s = fmt.Sprintf("endpoints: %d good", len(o.Good))

	if len(o.Bad) > 0 {
		s += fmt.Sprintf(", %d bad", len(o.Bad))
	}

	if len(o.Suspicious) > 0 {
		s += fmt.Sprintf(", %d suspicious", len(o.Suspicious))
	}

	if len(o.Unknown) > 0 {
		s += fmt.Sprintf(", %d unknown", len(o.Unknown))
	}

	return
}

func (o *EndpointOverview) GetStatus() int {
	if len(o.Bad) > 0 {
		return check.Critical
	} else if len(o.Suspicious) > 0 || len(o.Unknown) > 0 {
		return check.Warning
	} else {
		return check.OK
	}
}

func (o *EndpointOverview) GetOutput(limit int) (s string) {
	if o.Total == 0 {
		return
	}

	s = "\n## Endpoints\n"

	if len(o.Bad) > 0 {
		s += fmt.Sprintf("bad: %s\n", JoinEmphasis(o.Bad, ", ", limit))
	}

	if len(o.Suspicious) > 0 {
		s += fmt.Sprintf("suspicious: %s\n", JoinEmphasis(o.Suspicious, ", ", limit))
	}

	if len(o.Unknown) > 0 {
		s += fmt.Sprintf("unknown: %s\n", JoinEmphasis(o.Unknown, ", ", limit))
	}

	return
}

func (o *EndpointOverview) GetPerfdata() string {
	return PerfdataList{
		{Name: "endpoints_total", Value: fmt.Sprintf("%d", o.Total)},
		{Name: "endpoints_good", Value: fmt.Sprintf("%d", len(o.Good))},
		{Name: "endpoints_bad", Value: fmt.Sprintf("%d", len(o.Bad))},
		{Name: "endpoints_suspicious", Value: fmt.Sprintf("%d", len(o.Suspicious))},
		{Name: "endpoints_unknown", Value: fmt.Sprintf("%d", len(o.Unknown))},
	}.String()
}
