package config

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
	"time"
)

var GlobalConfiguration ReverseProxyConfiguration
var redisClient *redis.Client

type ServerConfiguration struct {
	Name       string `yaml:"name"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	UrlMapping string `yaml:"urlMapping"`
	Protocol   string `yaml:"protocol"`
}

type RedisConfiguration struct {
	Host                 string `yaml:"host"`
	Port                 string `yaml:"port"`
	SecondsToExpireCache string `yaml:"secondsToExpireCache"`
}

type ReverseProxyConfiguration struct {
	ServerConfigurations []ServerConfiguration `yaml:"servers"`
	RedisConfiguration   RedisConfiguration    `yaml:"redis"`
}

func LoadConfiguration() {
	configFile, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(configFile, &GlobalConfiguration); err != nil {
		log.Fatal(err)
	}
}

func (reverseProxyConfiguration ReverseProxyConfiguration) GetServerUrl(url string) string {
	for _, serverConfiguration := range reverseProxyConfiguration.ServerConfigurations {
		if strings.Contains(url, serverConfiguration.UrlMapping) {
			_, method, _ := strings.Cut(url, serverConfiguration.UrlMapping)
			return fmt.Sprintf("%s://%s:%s%s", serverConfiguration.Protocol, serverConfiguration.Host, serverConfiguration.Port, method)
		}
	}
	return ""
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     getRedisConnectionUrl(),
			Password: "",
			DB:       0,
		})
	}
	return redisClient
}

func getRedisConnectionUrl() string {
	return fmt.Sprintf("%s:%s", GlobalConfiguration.RedisConfiguration.Host, GlobalConfiguration.RedisConfiguration.Port)
}

func GetRedisSecondsToExpireCache() time.Duration {
	duration, err := time.ParseDuration(GlobalConfiguration.RedisConfiguration.SecondsToExpireCache + "s")
	if err != nil {
		return time.Second * 60
	}
	return duration
}
