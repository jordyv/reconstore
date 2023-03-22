package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Initialize() {
	viper.SetConfigName("reconstore")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")

	initDefaults()

	err := viper.ReadInConfig()
	if err != nil {
		logrus.WithError(err).Fatal("cannot read your config file")
	}
}

func initDefaults() {
	viper.SetDefault(DBType, "sqlite://reconstore.db")
}

func GetString(c string) string {
	return viper.GetString(c)
}
