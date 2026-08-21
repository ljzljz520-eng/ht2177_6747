package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	DatabasePath  string
	HTTPAddress   string
	Reviewer      string
	PolicyMinimum int
}

func Default() Config {
	return Config{DatabasePath: "store-ledger.db", HTTPAddress: "127.0.0.1:8090", Reviewer: "reviewer", PolicyMinimum: 70}
}

func Load() Config {
	c := Default()
	if value := strings.TrimSpace(os.Getenv("STORE_LEDGER_DB")); value != "" {
		c.DatabasePath = value
	}
	if value := strings.TrimSpace(os.Getenv("STORE_LEDGER_HTTP")); value != "" {
		c.HTTPAddress = value
	}
	if value := strings.TrimSpace(os.Getenv("STORE_LEDGER_REVIEWER")); value != "" {
		c.Reviewer = value
	}
	return c
}

func (c Config) Validate() error {
	if c.DatabasePath == "" {
		return errors.New("database path is required")
	}
	if c.HTTPAddress == "" {
		return errors.New("http address is required")
	}
	if c.Reviewer == "" {
		return errors.New("reviewer is required")
	}
	if c.PolicyMinimum < 0 || c.PolicyMinimum > 100 {
		return errors.New("policy minimum is invalid")
	}
	return nil
}

func (c Config) AbsoluteDatabasePath(base string) string {
	if filepath.IsAbs(c.DatabasePath) {
		return c.DatabasePath
	}
	return filepath.Join(base, c.DatabasePath)
}

func (c Config) WithDatabase(path string) Config     { c.DatabasePath = path; return c }
func (c Config) WithAddress(address string) Config   { c.HTTPAddress = address; return c }
func (c Config) WithReviewer(reviewer string) Config { c.Reviewer = reviewer; return c }
