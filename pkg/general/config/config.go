package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	env         string
	viperConfig *viper.Viper
}

func InitConfigs(environment string) *AppConfig {

	appConfig := &AppConfig{
		env:         environment,
		viperConfig: viper.New(),
	}
	return appConfig
}

func (a *AppConfig) LoadConfigs() {

	a.viperConfig.SetConfigName(fmt.Sprintf("config.%s", a.env))
	a.viperConfig.SetConfigType("yaml")
	a.viperConfig.AddConfigPath("../../configurations")

	err := a.viperConfig.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config files : %v", err)
	}

	a.viperConfig.SetConfigName(fmt.Sprintf("config.%s", "common"))
	err = a.viperConfig.MergeInConfig()
	if err != nil {
		log.Printf("Error reading config files : %v", err)
	}
}

func (a *AppConfig) PrintAllKeys() {
	settings := a.viperConfig.AllSettings()
	for key, value := range settings {
		log.Printf("%s: %v\n", key, value)
	}
}

func (a *AppConfig) GetConfig(key string) string {
	return a.viperConfig.GetString(key)
}

func (a *AppConfig) SetConfig(key, value string) {
	a.viperConfig.Set(key, value)
}
