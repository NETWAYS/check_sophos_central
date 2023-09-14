package main

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/NETWAYS/check_sophos_central/api"
	"github.com/NETWAYS/go-check"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// Matches a list of regular expressions against a string.
func matches(input string, regexToExclude []string) bool {
	for _, regex := range regexToExclude {
		re := regexp.MustCompile(regex)
		if re.MatchString(input) {
			return true
		}
	}

	return false
}

type Config struct {
	ClientID         string
	ClientSecret     string
	APIBaseURL       string
	ShowAll          bool
	PageSize         uint32
	ExcludeAlerts    []string
	ExcludeEndpoints []string
}

func BuildConfigFlags(fs *pflag.FlagSet) (config *Config) {
	config = &Config{}

	fs.StringVar(&config.ClientID, "client-id", "", "API Client ID (env:SOPHOS_CLIENT_ID)")
	fs.StringVar(&config.ClientSecret, "client-secret", "", "API Client Secret (env:SOPHOS_CLIENT_SECRET)")
	fs.BoolVar(&config.ShowAll, "show-all", false, "List all non-ok endpoints")
	fs.Uint32Var(&config.PageSize, "page-size", api.DefaultPageSize, "Amount of objects to fetch during each API call")
	fs.StringArrayVar(&config.ExcludeAlerts, "exclude-alert", []string{}, "Alerts to ignore. Can be used multiple times and supports regex.")          //nolint:lll
	fs.StringArrayVar(&config.ExcludeEndpoints, "exclude-endpoint", []string{}, "Endpoints to ignore. Can be used multiple times and supports regex.") //nolint:lll

	fs.StringVar(&config.APIBaseURL, "api", api.DefaultURL, "API Base URL")

	return
}

func (c *Config) SetFromEnv() {
	if c.ClientID == "" {
		c.ClientID = os.Getenv("SOPHOS_CLIENT_ID")
	}

	if c.ClientSecret == "" {
		c.ClientSecret = os.Getenv("SOPHOS_CLIENT_SECRET")
	}
}

func (c *Config) Validate() error {
	if c.ClientID == "" || c.ClientSecret == "" {
		return errors.New("client-id and client-secret are required")
	}

	return nil
}

func (c *Config) Run() (rc int, output string, err error) {
	// Setup API client.
	client := api.NewClient(c.ClientID, c.ClientSecret)
	client.PageSize = c.PageSize

	err = client.WhoAmI()
	if err != nil {
		return
	}

	log.WithField("context-id", client.UserInfo.ID).Debug("successfully authenticated with the API")

	// Retrieve and check endpoints.
	endpoints, names, err := CheckEndpoints(client, c.ExcludeEndpoints)
	if err != nil {
		return
	}

	// Retrieve and check alerts.
	alerts, err := CheckAlerts(client, names, c.ExcludeAlerts)
	if err != nil {
		return
	}

	// Build output.
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

	// nolint: gocritic
	if rcAlerts == check.Critical || rcEndpoints == check.Critical {
		rc = check.Critical
	} else if rcEndpoints > rcAlerts {
		rc = rcEndpoints
	} else {
		rc = rcAlerts
	}

	// Add Perfdata.
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
