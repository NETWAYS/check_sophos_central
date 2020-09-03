package api_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_GetAlerts(t *testing.T) {
	c := envClient(t)

	err := c.WhoAmI()
	assert.NoError(t, err)

	_, err = c.GetAlerts()
	assert.NoError(t, err)
}
