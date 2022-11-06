package test

import (
	"github.com/alicebob/miniredis"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"reverseproxy/src/config"
)

func LoadRedisMock() {
	redis := miniredis.NewMiniRedis()
	redis.StartAddr("localhost:6379")
}

func LoadReverseProxyConfiguration() {
	configFile, _ := os.ReadFile("../config_test.yml")

	if err := yaml.Unmarshal(configFile, &config.GlobalConfiguration); err != nil {
		log.Fatal(err)
	}
}

func StartBackend() {
	http.HandleFunc("/ping", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	}))
	go http.ListenAndServe(":9000", nil)
}
