package config

import "github.com/spf13/viper"

type DatabaseCfg struct {
	Host string
	Port int
	User string
	Pwd  string
}

type Config struct {
	Database DatabaseCfg
}

func NewConfig() Config {
	return Config{
		Database: DatabaseCfg{
			Host: viper.GetString("db.host"),
			Port: viper.GetInt("db.port"),
			User: viper.GetString("db.user"),
			Pwd:  viper.GetString("db.pwd"),
		},
	}
}
