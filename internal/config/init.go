package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
)

func Initialize(fileName string) {
	dir, fileName := path.Split(fileName)

	viper.SetConfigName(fileName)
	viper.SetConfigType("yaml")
	if dir != "" {
		viper.AddConfigPath(dir)
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME")

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
