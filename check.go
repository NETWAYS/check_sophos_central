package main

import (
	"errors"
	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"os"
	"strings"
)

type Config struct {
	ClientID     string
	ClientSecret string
	ApiBaseUrl   string
	ShowAll      bool
}

func BuildConfigFlags(fs *pflag.FlagSet) (config *Config) {
	config = &Config{}

	fs.StringVar(&config.ClientID, "client-id", "", "API Client ID (env:SOPHOS_CLIENT_ID)")
	fs.StringVar(&config.ClientSecret, "client-secret", "", "API Client Secret (env:SOPHOS_CLIENT_SECRET)")
	fs.BoolVar(&config.ShowAll, "show-all", false, "List all non-ok endpoints")

	fs.StringVar(&config.ApiBaseUrl, "api", api.DefaultURL, "API Base URL")

	return
}

func (c *Config) SetFromEnv() {
	if c.ClientID == "" {
		c.ClientID = os.Getenv("SOPHOS_CLIENT_ID")
	}

	if c.ClientSecret == "" {
		c.ClientSecret = os.Getenv("SOPHOS_CLIENT_SECRET")
	}

	return
}

func (c *Config) Validate() (err error) {
	if c.ClientID == "" || c.ClientSecret == "" {
		err = errors.New("client-id and client-secret are required")
		return
	}

	return
}

func (c *Config) Run() (rc int, output string, err error) {
	// Setup API client
	client := api.NewClient(c.ClientID, c.ClientSecret)

	err = client.WhoAmI()
	if err != nil {
		return
	}

	log.WithField("context-id", client.UserInfo.ID).Debug("successfully authenticated with the API")

	// Retrieve and check alerts
	alerts, err := CheckAlerts(client)
	if err != nil {
		return
	}

	// Retrieve and check endpoints
	endpoints, err := CheckEndpoints(client)
	if err != nil {
		return
	}

	// Build output
	limit := 5
	if c.ShowAll {
		limit = 0
	}

	output = alerts.GetSummary() + " - " + endpoints.GetSummary() + "\n"
	output += alerts.GetOutput()
	output += endpoints.GetOutput(limit)

	// Build rc
	rcAlerts := alerts.GetStatus()
	rcEndpoints := endpoints.GetStatus()

	if rcAlerts == check.Critical || rcEndpoints == check.Critical {
		rc = check.Critical
	} else if rcEndpoints > rcAlerts {
		rc = rcEndpoints
	} else {
		rc = rcAlerts
	}

	// Add Perfdata
	output += "| " + alerts.GetPerfdata() + " " + endpoints.GetPerfdata()

	return
}

func JoinEmphasis(elems []string, sep string, limit int) string {
	if limit > 0 && len(elems) > limit {
		elems = elems[0:limit]
		elems = append(elems, "...")
	}

	return strings.Join(elems, sep)
}
