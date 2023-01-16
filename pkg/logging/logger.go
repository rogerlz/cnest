package logging

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
)

func NewLogger(logLevel string) log.Logger {

	var (
		logger log.Logger
		lvl    level.Option
	)

	switch logLevel {
	case "quiet":
		lvl = level.AllowError()
	case "verbose":
		lvl = level.AllowInfo()
	case "debug":
		lvl = level.AllowDebug()
	default:
		panic("unexpected log level")
	}

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	logger = level.NewFilter(logger, lvl)

	return log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)
}
