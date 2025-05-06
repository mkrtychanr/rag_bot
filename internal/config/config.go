package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   Logger    `mapstructure:"logger" yaml:"logger"`
	Profile  Profiling `mapstructure:"profiling" yaml:"profiling"`
	Bot      Bot       `mapstructure:"bot" yaml:"bot"`
	Postrges Postrges  `mapstructure:"postgres" yaml:"postgres"`
	Rag      Rag       `mapstructure:"rag" yaml:"rag"`
}

func (c *Config) Validate() error {
	// if err := c.Cache.Validate(); err != nil {
	// 	return fmt.Errorf("cache config validation failed. %w", err)
	// }

	// if err := c.Chain.Validate(); err != nil {
	// 	return fmt.Errorf("chain config validation failed. %w", err)
	// }

	// return c.Admin.Validate()
	return nil
}

func setupDefaultConfigPaths() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/rb")

	workDir, err := os.Getwd()
	if err == nil {
		env := os.Getenv("RB_ENV")

		if env != "" {
			viper.AddConfigPath(path.Join(workDir, "config", env))
		}

		viper.AddConfigPath(path.Join(workDir, "config"))
	}
}

// FindConfig is searching config file specified with path.
//
// If there is no such file config will be created using env variables.
func FindConfig(p string) (*Config, error) {
	cfg, err := readCofigFromFile(p)
	if err != nil {
		return nil, fmt.Errorf("could not find config. %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("loaded configuration is invalid. %w", err)
	}

	return cfg, nil
}

func readCofigFromFile(p string) (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("RB")
	viper.AutomaticEnv()

	if p == "" {
		setupDefaultConfigPaths()
	} else {
		absP, err := filepath.Abs(p)
		if err != nil {
			return nil, fmt.Errorf("could not convert path to absolute. %w", err)
		}

		viper.SetConfigFile(absP)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("can't find config. %w", err)
	}

	var cfg Config

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("could not unmarshal configuration. %w", err)
	}

	return &cfg, nil
}
