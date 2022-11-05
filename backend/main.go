package main

import (
	"backend/src/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ping", handler.Ping)
	http.HandleFunc("/sum", handler.Sum)

	log.Fatal(http.ListenAndServe(":9000", nil))
}
