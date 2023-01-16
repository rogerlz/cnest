package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/rogerlz/cnest/pkg/config"
	"github.com/rogerlz/cnest/pkg/executor"
	"github.com/rogerlz/cnest/pkg/logging"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
)

const (
	appName = "cnest"
)

func main() {
	v, p := viper.New(), pflag.NewFlagSet(appName, pflag.ExitOnError)
	config, configFileNotFound := config.Configure(v, p)
	logger := logging.NewLogger(config.Log)

	if configFileNotFound {
		level.Debug(logger).Log("msg", "configuration file not found, using defaults")
	}

	level.Info(logger).Log("msg", "starting application")

	var g run.Group

	// Start the executor module
	executor := executor.New(logger, config.Executor)
	{
		g.Add(func() error {
			return executor.RunCmd()
		}, func(error) {
			executor.Stop()
		})
	}

	// Listen for termination signals.
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			return interrupt(logger, cancel)
		}, func(error) {
			close(cancel)
		})
	}

	if err := g.Run(); err != nil {
		level.Error(logger).Log("err", fmt.Sprintf("%+v", errors.Wrapf(err, "failed to start")))
		os.Exit(1)
	}

	level.Info(logger).Log("msg", "exiting")
}

func interrupt(logger log.Logger, cancel <-chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-c:
		level.Info(logger).Log("msg", "caught signal. Exiting.", "signal", s)
		return nil
	case <-cancel:
		return errors.New("canceled")
	}
}
