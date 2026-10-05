// Package config loads and validates the server configuration.
package config

import (
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
	cfg := Config{Addr: getenv("HASHPLACE_ADDR")}
	if cfg.Addr == "" {
		cfg.Addr = defaultAddr
	}

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

	n, err := strconv.Atoi(port)
	if err != nil || n < 0 || n > 65535 {
		return fmt.Errorf("port %q must be a number between 0 and 65535", port)
	}

	return nil
}
