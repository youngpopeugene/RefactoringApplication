package config

import (
	"github.com/spf13/viper"
)

func init() {
	bindEnvs()
	setDefaultValues()
}

func failOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func bindEnvs() {
	viper.AutomaticEnv()

	failOnError(viper.BindEnv("db.host", "DB_HOST"))
	failOnError(viper.BindEnv("db.port", "DB_PORT"))
	failOnError(viper.BindEnv("db.user", "DB_USER"))
	failOnError(viper.BindEnv("db.pwd", "DB_PWD"))
}

func setDefaultValues() {
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 5432)
	viper.SetDefault("db.user", "guest")
	viper.SetDefault("db.pwd", "guest")
}
