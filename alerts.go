package main

import (
	"fmt"
	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/NETWAYS/go-check"
	"strings"
)

type AlertOverview struct {
	Total  int
	High   int
	Medium int
	Low    int
	Output []string
}

func CheckAlerts(client *api.Client) (o *AlertOverview, err error) {
	o = &AlertOverview{}

	alerts, err := client.GetAlerts()
	if err != nil {
		return
	}

	for _, alert := range alerts {
		o.Total++

		switch strings.ToLower(alert.Severity) {
		case "high":
			o.High++
		case "medium":
			o.Medium++
		default:
			o.Low++
		}

		output := fmt.Sprintf("[%s] %s: %s group=%s",
			alert.RaisedAt, alert.Type, alert.Description, alert.GroupKey)
		o.Output = append(o.Output, output)
	}

	return
}

func (o *AlertOverview) GetSummary() string {
	if o.Total == 0 {
		return "no alerts"
	}

	var states []string
	if o.High > 0 {
		states = append(states, fmt.Sprintf("%d high", o.High))
	}

	if o.Medium > 0 {
		states = append(states, fmt.Sprintf("%d medium", o.Medium))
	}

	if o.Low > 0 {
		states = append(states, fmt.Sprintf("%d low", o.Low))
	}

	return "alerts: " + strings.Join(states, ", ")
}

func (o *AlertOverview) GetStatus() int {
	if o.High > 0 {
		return check.Critical
	} else if o.Medium > 0 || o.Low > 0 {
		return check.Warning
	} else {
		return check.OK
	}
}

func (o *AlertOverview) GetOutput() (s string) {
	if o.Total == 0 {
		return
	}

	s = "\n## Alerts"
	s += strings.Join(o.Output, "\n")
	s += "\n"

	return
}
