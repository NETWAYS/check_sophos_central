package api_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClient_WhoAmI(t *testing.T) {
	c := envClient(t)

	err := c.WhoAmI()
	assert.NoError(t, err)
	assert.NotNil(t, c.UserInfo)
	assert.NotEmpty(t, c.UserInfo.ID)
	assert.NotEmpty(t, c.UserInfo.IDType)
	assert.NotEmpty(t, c.DataURL)
}
