package logging

import "errors"

type Config struct {
	Format string
	Level  string
}

func (c Config) Validate() error {
	if c.Format != "json" && c.Format != "logfmt" {
		return errors.New("log level must be json or logfmt")
	}

	if c.Level != "error" && c.Level != "warn" && c.Level != "info" && c.Level != "debug" {
		return errors.New("log level must be error, warn, info or debug")
	}

	return nil
}
