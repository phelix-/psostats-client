// psostats client config
package config

import (
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	defaultUIRefreshRate = time.Second / 5
)

type Config struct {
	ServerBaseUrl *string `yaml:"serverBaseUrl"`
	UiFps         *int    `yaml:"uiFps"`
	User          *string `yaml:"user"`
	Password      *string `yaml:"password"`
	AutoUpload    *bool   `yaml:"autoUpload"`
}

func (config *Config) GetUiRefreshRate() time.Duration {
	uiRefreshRate := defaultUIRefreshRate
	if config.UiFps != nil {
		if *config.UiFps <= 0 || *config.UiFps > 30 {
			log.Printf("uiFps must be between 0 and 30 but was %v, falling back to default(5)", *config.UiFps)
		} else {
			uiRefreshRate = time.Second / time.Duration(*config.UiFps)
		}
	}
	return uiRefreshRate
}

func (config *Config) GetServerBaseUrl() string {
	if config.ServerBaseUrl != nil {
		return *config.ServerBaseUrl
	} else {
		return "https://psostats.com"
	}
}

func ReadFromFile(fileLocation string) *Config {
	config := Config{}
	data, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Fatalf("Error reading configuration file %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing config file %v", err)
	}
	return &config
}

func (config *Config) AutoUploadEnabled() bool {
	return config.AutoUpload == nil || *config.AutoUpload
}
