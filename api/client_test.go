package api_test

import (
	"github.com/NETWAYS/check_sophos_central/api"
	"os"
	"testing"
)

func envClient(t *testing.T) *api.Client {
	id := os.Getenv("SOPHOS_CLIENT_ID")
	secret := os.Getenv("SOPHOS_CLIENT_SECRET")

	if id == "" || secret == "" {
		t.Skip("SOPHOS_CLIENT_ID and SOPHOS_CLIENT_SECRET must be set!")
	}

	return api.NewClient(id, secret)
}
