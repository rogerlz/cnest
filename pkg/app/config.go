package app

import "errors"

const (
	appName = "cnest"
)

type Config struct {
	LogLevel string `ini:"log_level"`
}

func (c Config) Validate() error {
	if c.LogLevel != "quiet" && c.LogLevel != "verbose" && c.LogLevel != "debug" {
		return errors.New("log level must be quiet, verbose or debug")
	}

	return nil
}
