package executor

import (
	"errors"
	"os"
)

type Config struct {
	CmdPath string
	CmdArgs string
}

func (c Config) Validate() error {
	if _, err := os.Stat(c.CmdPath); err != nil {
		return errors.New("command path does not exist")
	}

	return nil
}
