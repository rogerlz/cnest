package logging

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func NewLogger(config Config) log.Logger {

	var (
		logger log.Logger
		lvl    level.Option
	)

	switch config.Level {
	case "error":
		lvl = level.AllowError()
	case "warn":
		lvl = level.AllowWarn()
	case "info":
		lvl = level.AllowInfo()
	case "debug":
		lvl = level.AllowDebug()
	default:
		panic("unexpected log level")
	}

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	if config.Format == "json" {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	}

	logger = level.NewFilter(logger, lvl)

	return log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
}
