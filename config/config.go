package config

import (
	"strings"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("MINER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()
	viper.SetDefault("connectionStrings.redis", "localhost:6379")
	viper.SetDefault("connectionStrings.mongodb", "mongodb://localhost:27017")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
}
