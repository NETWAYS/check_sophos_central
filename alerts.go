package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/NETWAYS/go-check"
)

type AlertOverview struct {
	Total  int
	High   int
	Medium int
	Low    int
	Output []string
}

// Retrieve and process Alerts.
// alertsToExclude is a list of strings that can the used to exclude alerts.
func CheckAlerts(client *api.Client, names EndpointNames, alertsToExclude []string) (o *AlertOverview, err error) {
	o = &AlertOverview{}

	alerts, err := client.GetAlerts()
	if err != nil {
		return
	}

	for _, alert := range alerts {
		if matches(alert.String(), alertsToExclude) {
			// If the alert matches a regex from the list we can skip it.
			continue
		}

		o.Total++

		switch strings.ToLower(alert.Severity) {
		case "high":
			o.High++
		case "medium":
			o.Medium++
		default:
			o.Low++
		}

		agentName := alert.ManagedAgent.ID
		if val, ok := names[agentName]; ok {
			agentName = val
		}

		output := fmt.Sprintf("%s [%s] %s (%s) %s",
			//nolint: gosmopolitan
			alert.RaisedAt.Local().Format("2006-01-02 15:04 MST"),
			alert.Severity, agentName, alert.Product, alert.Description)
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

// nolint: gocritic
func (o *AlertOverview) GetStatus() int {
	if o.High > 0 {
		return check.Critical
	} else if o.Medium > 0 || o.Low > 0 {
		return check.Warning
	}

	return check.OK
}

func (o *AlertOverview) GetOutput() (s string) {
	if o.Total == 0 {
		return
	}

	s = "\n## Alerts\n"
	s += strings.Join(o.Output, "\n")
	s += "\n"

	return
}

func (o *AlertOverview) GetPerfdata() string {
	return PerfdataList{
		{Name: "alerts", Value: strconv.Itoa(o.Total)},
		{Name: "alerts_high", Value: strconv.Itoa(o.High)},
		{Name: "alerts_medium", Value: strconv.Itoa(o.Medium)},
		{Name: "alerts_low", Value: strconv.Itoa(o.Low)},
	}.String()
}
