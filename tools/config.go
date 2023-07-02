package tools

import (
	"github.com/spf13/viper"
)

type Config struct {
	ReindexerServerAddress string `mapstructure:"REINDEXER_SERVER_ADDRESS"`
	DatabaseName           string `mapstructure:"DATABASE_NAME"`
	Namespace              string `mapstructure:"NAMESPACE"`
	AppServerAddress       string `mapstructure:"APP_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath("configs")

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
