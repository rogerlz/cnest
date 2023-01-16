package camera

import (
	"github.com/go-kit/kit/log"

	"github.com/rogerlz/cnest/pkg/executor"
)

type Camera struct {
	Config   Config
	Executor executor.Executor
}

func New(l log.Logger, c Config) *Camera {

	var (
		cmdPath string
		cmdArgs []string
	)

	switch c.Mode {
	case "mjpg":
		cmdPath, cmdArgs = GetMjpgArguments(c)
	case "rtsp":
		cmdPath, cmdArgs = GetRtspArguments(c)
	}

	return &Camera{
		Config:   c,
		Executor: *executor.New(l, cmdPath, cmdArgs),
	}
}

func (c *Camera) Run() error {
	return c.Executor.RunCmd()
}

func (c *Camera) Stop() {
	c.Executor.Stop()
}

/*

const (
	argsRegex = `[^\s"]+|"([^"]*)"`
)

r := regexp.MustCompile(argsRegex)
parsedArgs := r.FindAllString(args, -1)

*/
