package config

import (
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var (
	Port       int
	Name       string
	Mode       string
	Level      string
	Output     string
	LogFile    string
	RetryCount int
)

type Config struct {
	Port       int    `mapstructure:"port"`
	Name       string `mapstructure:"name"`
	Mode       string `mapstructure:"mode"`
	Level      string `mapstructure:"level"`
	Output     string `mapstructure:"output"` // "stdout" 或 "file"
	LogFile    string `mapstructure:"log_file"`
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

	// v.SetEnvPrefix("OCTO")                             //只认OCTO开头的环境变量
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) //将 "." 替换为 "_"
	v.AutomaticEnv()
	if err := v.Unmarshal(&GlobalConfig); err != nil {
		return err
	}
	Port = GlobalConfig.Port
	Name = GlobalConfig.Name
	Mode = GlobalConfig.Mode
	Level = GlobalConfig.Level
	Output = GlobalConfig.Output
	LogFile = GlobalConfig.LogFile
	RetryCount = GlobalConfig.RetryCount
	// fmt.Printf("%+v", GlobalConfig)
	return nil
}
