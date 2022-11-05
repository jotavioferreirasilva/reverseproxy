package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

var Servers Server

type ServerConfiguration struct {
	ServerName string `yaml:"name"`
	ServerHost string `yaml:"host"`
	ServerPort string `yaml:"port"`
	ServerUrl  string `yaml:"urlMapping"`
}

type Server struct {
	ServerConfigurations []ServerConfiguration `yaml:"servers"`
}

func LoadConfiguration() {
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(configFile, &Servers); err != nil {
		log.Fatal(err)
	}
}

func (server Server) GetServerUrl(url string) string {
	for _, serverConfiguration := range server.ServerConfigurations {
		if strings.Contains(url, serverConfiguration.ServerUrl) {
			_, after, _ := strings.Cut(url, serverConfiguration.ServerUrl)
			return fmt.Sprintf("http://%s:%s%s", serverConfiguration.ServerHost, serverConfiguration.ServerPort, after)
		}
	}
	return ""
}
