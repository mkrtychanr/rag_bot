package config

import "fmt"

type Postrges struct {
	UserName     string `mapstructure:"userName" yaml:"userName"`
	Password     string `mapstructure:"password" yaml:"password"`
	Host         string `mapstructure:"host" yaml:"host"`
	Port         string `mapstructure:"port" yaml:"port"`
	DataBaseName string `mapstructure:"dataBaseName" yaml:"dataBaseName"`
	SSL          bool   `mapstructure:"ssl" yaml:"ssl"`
}

func (p Postrges) BuildPostgresConnectionString() string {
	ssl := "disable"
	if p.SSL {
		ssl = "enable"
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", p.UserName, p.Password, p.Host, p.Port, p.DataBaseName, ssl)
}
