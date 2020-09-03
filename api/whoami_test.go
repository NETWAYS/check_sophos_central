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

func TestClient_WhoAmI_Mock(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	err := c.WhoAmI()
	assert.NoError(t, err)
	assert.NotNil(t, c.UserInfo)
	assert.Equal(t, "57ca9a6b-885f-4e36-95ec-290548c26059", c.UserInfo.ID)
	assert.Equal(t, "tenant", c.UserInfo.IDType)
	assert.Equal(t, "https://api-eu01.central.sophos.com", c.DataURL)
}
