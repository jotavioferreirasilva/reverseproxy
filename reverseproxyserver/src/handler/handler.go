package handler

import (
	"fmt"
	"io"
	"net/http"
	"reverseproxy/src/config"
)

func ReverseProxyHandler(writer http.ResponseWriter, request *http.Request) {
	client := &http.Client{}

	var url = fmt.Sprintf("%s", config.Servers.GetServerUrl(request.URL.RequestURI()))

	requestToBackend, err := http.NewRequest(request.Method, url, request.Body)
	requestToBackend.Header = request.Header

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := client.Do(requestToBackend)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(response.Body)

	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(response.StatusCode)
	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
