package config

import (
	"log"

	"github.com/spf13/viper"
)

var Config ConfigType

type ConfigType struct {
	Env      string `mapstructure:"env"`
	Port     int    `mapstructure:"port"`
	Debug    bool   `mapstructure:"debug"`
	Secret   string `mapstructure:"string"`
	Database struct {
		URL        string `mapstructure:"url"`
		SqlDir     string `mapstructure:"sqldir"`
		Migrations string `mapstructure:"migrations"`
	} `mapstructure:"database"`
	Sendgrid struct {
		URL    string `mapstructure:"url"`
		ApiKey string `mapstructure:"api_key"`
	} `mapstructure:"sendgrid"`
	Wellet struct {
		Agent       string `mapstructure:"agent"`
		AgentApiKey string `mapstructure:"agent_api_key"`
		Connect     string `mapstructure:"connect"`
	} `mapstructure:"wallet"`
}

func Init(configPath string) {
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found: %s", err)
		} else {
			log.Fatalf("Error reading config file: %s", err)
		}
	}

	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatal(err)
	}

	log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
}
