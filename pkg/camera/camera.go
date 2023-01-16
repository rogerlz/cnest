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
	return &Camera{
		Config:   c,
		Executor: *executor.New(l, "sleep", "10"),
	}
}

func (c *Camera) Run() error {
	return c.Executor.RunCmd()
}

func (c *Camera) Stop() {
	c.Executor.Stop()
}
