package config

import (
	"../logger"
	"github.com/spf13/viper"
)

var (
	G_host    string = "localhost:1234"
	G_wav_dir string = "wav"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("ini")

	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error.Println("Fatal error config file: " + err.Error())
		logger.Info.Println("Using default setting")
	} else {
		if viper.GetString("default.host") != "" {
			G_host = viper.GetString("default.host")
		}
		if viper.GetString("default.wav_dir") != "" {
			G_wav_dir = viper.GetString("default.wav_dir")
		}
	}
}
