// Package config loads and validates the server configuration.
package config

import (
	"cmp"
	"fmt"
	"net"
	"strconv"
)

const defaultAddr = ":8080"

// Config holds the server settings.
type Config struct {
	// Addr is the TCP address the HTTP server listens on, for example ":8080".
	Addr string
}

// Load reads the configuration from getenv (usually os.Getenv) and validates it.
func Load(getenv func(string) string) (Config, error) {
	cfg := Config{Addr: cmp.Or(getenv("HASHPLACE_ADDR"), defaultAddr)}
	if err := validateAddr(cfg.Addr); err != nil {
		return Config{}, fmt.Errorf("invalid HASHPLACE_ADDR %q: %w", cfg.Addr, err)
	}

	return cfg, nil
}

func validateAddr(addr string) error {
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	_, err = strconv.ParseUint(port, 10, 16)
	return err
}
