package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/rogerlz/cnest/pkg/app"
	"github.com/rogerlz/cnest/pkg/camera"
	"github.com/rogerlz/cnest/pkg/logging"
	"github.com/spf13/pflag"
	"gopkg.in/ini.v1"

	"github.com/pkg/errors"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/run"
)

const (
	appName = "cnest"
)

func main() {
	p := pflag.NewFlagSet(appName, pflag.ExitOnError)

	configFile := p.StringP("config", "c", "", "configuration file")

	_ = p.Parse(os.Args[1:])

	if _, err := os.Stat(*configFile); err != nil {
		fmt.Printf("config not found: %v\n", err)
		os.Exit(1)
	}

	cfg, err := ini.Load(*configFile)
	if err != nil {
		fmt.Printf("failed to parse config file: %v\n", err)
		os.Exit(1)
	}

	appConfig := new(app.Config)
	cfg.Section("crowsnest").MapTo(appConfig)

	if err := appConfig.Validate(); err != nil {
		fmt.Printf("config validation: %v\n", err)
		os.Exit(1)
	}

	logger := logging.NewLogger(appConfig.LogLevel)

	level.Info(logger).Log("msg", "starting application")

	sections := cfg.SectionStrings()

	var g run.Group

	// find cameras config
	camerasFound := []string{}
	for _, s := range sections {
		found, _ := regexp.MatchString(`cam\s[0-9]`, s)
		if found {
			camerasFound = append(camerasFound, s)
		}
	}
	// TODO log cameras found here

	// Loop cameras found
	for _, c := range camerasFound {

		// Load and parse camera config
		cameraConfig := new(camera.Config)
		cfg.Section(c).MapTo(cameraConfig)

		// Validate camera config
		if err := cameraConfig.Validate(); err != nil {
			fmt.Printf("camera config validation: %v\n", err)
			os.Exit(1)
		}

		// Start a camera thread
		newCamera := camera.New(logger, *cameraConfig)
		{
			g.Add(func() error {
				return newCamera.Run()
			}, func(error) {
				newCamera.Stop()
			})
		}
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
