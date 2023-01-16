package config

import (
	"github.com/pkg/errors"
	"github.com/rogerlz/cnest/pkg/executor"
	"github.com/rogerlz/cnest/pkg/logging"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Log      logging.Config
	Executor executor.Config
}

func (c Config) Validate() error {
	if err := c.Log.Validate(); err != nil {
		return err
	}

	if err := c.Executor.Validate(); err != nil {
		return err
	}

	return nil
}

func Configure(v *viper.Viper, p *pflag.FlagSet) (Config, bool) {

	// logging defaults
	v.SetDefault("log.format", "logfmt")
	v.SetDefault("log.level", "info")

	// executor defaults
	v.SetDefault("executor.cmdPath", "/bin/echo")
	v.SetDefault("executor.cmdArgs", "argument")

	v.AddConfigPath(".")
	v.AddConfigPath("/etc/cnest")
	v.SetConfigName("config")
	v.SetConfigType("toml")

	err := v.ReadInConfig()
	_, configFileNotFound := err.(viper.ConfigFileNotFoundError)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		panic(errors.Wrap(err, "failed to unmarshal configuration"))
	}

	if err := config.Validate(); err != nil {
		panic(errors.Wrap(err, "failed to validate configuration"))
	}

	return config, configFileNotFound
}
