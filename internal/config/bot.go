package config

type Bot struct {
	Token   string `mapstructure:"token" yaml:"token"`
	Offset  int    `mapstructure:"offset" yaml:"offset"`
	Timeout int    `mapstructure:"timeout" yaml:"timeout"`
}
