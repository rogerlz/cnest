package executor

import (
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/pkg/errors"
)

const (
	argsRegex = `[^\s"]+|"([^"]*)"`
)

type Executor struct {
	CmdPath string
	CmdArgs []string

	logger log.Logger
	cmd    *exec.Cmd
}

func New(l log.Logger, path, args string) *Executor {
	// TODO validate if path exists here

	r := regexp.MustCompile(argsRegex)
	parsedArgs := r.FindAllString(args, -1)
	cmd := exec.Command(path, parsedArgs[0:]...)

	return &Executor{
		CmdPath: path,
		CmdArgs: parsedArgs,
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
