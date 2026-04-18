package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port       int    `mapstructure:"port"`
	Name       string `mapstructure:"name"`
	Mode       string `mapstructure:"mode"`
	Level      string `mapstructure:"level"`
	Output     string `mapstructure:"output"` // "stdout" 或 "file"
	FileDir    string `mapstructure:"filedir"`
	RetryCount int    `mapstructure:"retry_count"`
}

var GlobalConfig Config

func LoadConfig(path string) error {

	if err := godotenv.Load(); err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return err
		}
	}

	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return err
	}

	v.SetEnvPrefix("OCTO") //将 "." 替换为 "_"
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	if err := v.Unmarshal(&GlobalConfig); err != nil {
		return err
	}
	return nil
}
