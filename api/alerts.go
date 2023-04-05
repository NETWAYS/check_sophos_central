package api

import (
	"encoding/json"
	"fmt"
	"time"
)

type Alert struct {
	ID string `json:"id"`
	// Alert categories.
	//
	// The following values are allowed:
	// azure, adSync, applicationControl, appReputation, blockListed, connectivity, cwg, denc, downloadReputation,
	// endpointFirewall, fenc, forensicSnapshot, general, iaas, iaasAzure, isolation, malware, mtr, mobiles, policy,
	// protection, pua, runtimeDetections, security, smc, systemHealth, uav, uncategorized, updating, utm, virt,
	// wireless, xgEmail
	Category     string            `json:"category"`
	Description  string            `json:"description"`
	GroupKey     string            `json:"groupKey"`
	ManagedAgent AlertManagedAgent `json:"managedAgent"`
	// Product types.
	//
	// The following values are allowed:
	// other, endpoint, server, mobile, encryption, emailGateway, webGateway, phishThreat, wireless, iaas, firewall
	Product string `json:"product"`
	// When the alert was triggered.
	RaisedAt time.Time `json:"raisedAt"`
	// Severity levels for alerts.
	//
	// The following values are allowed:
	// high, medium, low
	Severity string `json:"severity"`
	// Alert type.
	Type string `json:"type"`
}

func (a *Alert) String() string {
	return fmt.Sprintf(
		"%v [%v] %v %v %v %v %v %v",
		a.RaisedAt.Format("2006-01-02 15:04"),
		a.Severity,
		a.ID,
		a.Product,
		a.Description,
		a.Type,
		a.Category,
		a.GroupKey,
	)
}

type AlertManagedAgent struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func (c *Client) GetAlerts() (alerts []*Alert, err error) {
	req, err := c.NewDataRequest("GET", "common/v1/alerts", nil)
	if err != nil {
		return
	}

	items, err := c.GetResults(req)
	if err != nil {
		return
	}

	// retrieve items from response
	for _, item := range items {
		a := &Alert{}

		err = json.Unmarshal(item, a)
		if err != nil {
			return
		}

		alerts = append(alerts, a)
	}

	return
}
