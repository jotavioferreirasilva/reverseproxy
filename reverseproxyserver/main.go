package main

import (
	"log"
	"net/http"
	"reverseproxy/src/config"
	"reverseproxy/src/handler"
)

func main() {

	config.LoadConfiguration()

	http.HandleFunc("/", handler.ReverseProxyHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
