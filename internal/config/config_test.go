package config

import "testing"

func TestConfigDefaultsAndValidation(t *testing.T) {
	c := Default()
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	if c.WithDatabase("x.db").DatabasePath != "x.db" {
		t.Fatal("database override failed")
	}
}
