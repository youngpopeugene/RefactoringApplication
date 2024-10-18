package config

import "github.com/spf13/viper"

type PostgresCfg struct {
	Host string
	Port int
	User string
	Pwd  string
}

type Config struct {
	Postgres PostgresCfg
}

func NewConfig() Config {
	return Config{
		Postgres: PostgresCfg{
			Host: viper.GetString("db.host"),
			Port: viper.GetInt("db.port"),
			User: viper.GetString("db.user"),
			Pwd:  viper.GetString("db.pwd"),
		},
	}
}
