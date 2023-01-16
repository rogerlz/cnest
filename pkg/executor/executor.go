package executor

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
)

const (
	argsRegex = `[^\s"]+|"([^"]*)"`
)

type Executor struct {
	logger log.Logger
	config Config
	cmd    *exec.Cmd
}

func New(l log.Logger, c Config) *Executor {
	r := regexp.MustCompile(argsRegex)
	args := r.FindAllString(c.CmdArgs, -1)
	cmd := exec.Command(c.CmdPath, args[0:]...)

	return &Executor{
		logger: l,
		config: c,
		cmd:    cmd,
	}
}

func (e *Executor) RunCmd() error {
	level.Debug(e.logger).Log("msg", "running command", "path", e.config.CmdPath, "args", e.config.CmdArgs)

	e.cmd.Stderr = os.Stderr
	e.cmd.Stdout = os.Stdout

	if err := e.cmd.Start(); err != nil {
		return errors.New("error executing command")
	}

	if err := e.cmd.Wait(); err != nil {
		return errors.New("error executing command")
	}

	return nil
}

func (e *Executor) Stop() {
	level.Debug(e.logger).Log("msg", "stopping command")

	e.cmd.Stderr = nil
	e.cmd.Stdout = nil

	e.cmd.Process.Kill()
}
