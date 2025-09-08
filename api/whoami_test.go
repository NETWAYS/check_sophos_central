package api_test

import (
	"testing"
)

func TestClient_WhoAmI(t *testing.T) {
	c := envClient(t)

	err := c.WhoAmI()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c.UserInfo == nil {
		t.Fatalf("expected not nil")
	}
}

func TestClient_WhoAmI_Mock(t *testing.T) {
	c, cleanup := testClient()
	defer cleanup()

	err := c.WhoAmI()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if c.UserInfo == nil {
		t.Fatalf("expected not nil")
	}

	if "57ca9a6b-885f-4e36-95ec-290548c26059" != c.UserInfo.ID {
		t.Fatalf("expected %v, got %v", "57ca9a6b-885f-4e36-95ec-290548c26059", c.UserInfo.ID)
	}

	if "tenant" != c.UserInfo.IDType {
		t.Fatalf("expected %v, got %v", "tentant", c.UserInfo.IDType)
	}

	if "https://api-eu01.central.sophos.com" != c.DataURL {
		t.Fatalf("expected %v, got %v", "https://api-eu01.central.sophos.com", c.DataURL)
	}
}
