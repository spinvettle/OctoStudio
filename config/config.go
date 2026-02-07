package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
	Mode string `mapstructure:"mode"`
}
type LoggingConfig struct {
	Level   string `mapstructure:"level"`
	Output  string `mapstructure:"output"`  // "stdout" 或 "file"
	FileDir string `mapstructure:"filedir"` // 日志文件目录
}

type Config struct {
	ServerConfig  ServerConfig  `mapstructure:"server"`
	LoggingConfig LoggingConfig `mapstructure:"logging"`
}

var GlobalConfig Config

func LoadConfig() error {

	if err := godotenv.Load(); err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			return err
		}
	}

	v := viper.New()
	v.SetConfigFile("./config.yaml")
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
