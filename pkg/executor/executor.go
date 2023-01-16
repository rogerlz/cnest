package executor

import (
	"os"
	"os/exec"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
)

type Executor struct {
	CmdPath string
	CmdArgs []string

	logger log.Logger
	cmd    *exec.Cmd
}

func New(l log.Logger, cmdPath string, cmdArgs []string) *Executor {

	if _, err := os.Stat(cmdPath); err != nil {
		level.Error(l).Log("msg", "command to execute not found", "path", cmdPath, "args", strings.Join(cmdArgs, " "))
		os.Exit(1)
	}

	cmd := exec.Command(cmdPath, cmdArgs...)

	return &Executor{
		CmdPath: cmdPath,
		CmdArgs: cmdArgs,
		logger:  l,
		cmd:     cmd,
	}
}

func (e *Executor) RunCmd() error {
	level.Debug(e.logger).Log("msg", "running command", "path", e.CmdPath, "args", strings.Join(e.CmdArgs, " "))

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
