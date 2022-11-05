package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reverseproxy/src/config"
)

type ResponseCache struct {
	StatusCode int
	Body       []byte
	Present    bool
}

func ReverseProxyHandler(writer http.ResponseWriter, request *http.Request) {
	var url = config.GlobalConfiguration.GetServerUrl(request.URL.RequestURI())

	if url == "" {
		http.Error(writer, "", http.StatusNotFound)
		return
	}

	if request.Method == http.MethodGet && request.Header.Get("Cache-Control") != "no-store" {
		cachedResponse, err := getCachedResponse(url)
		if err != nil {
			fmt.Printf("error reading cache for %s. Error: %s\n", url, err.Error())
		}

		if cachedResponse.Present {
			writer.Header().Add("Cached", "true")
			writer.WriteHeader(cachedResponse.StatusCode)
			_, err := writer.Write(cachedResponse.Body)
			if err != nil {
				http.Error(writer, "error writing response", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	requestToBackend, err := http.NewRequest(request.Method, url, request.Body)
	if err != nil {
		http.Error(writer, "error creating request", http.StatusInternalServerError)
		return
	}
	requestToBackend.Header = request.Header

	client := &http.Client{}
	response, err := client.Do(requestToBackend)
	if err != nil {
		http.Error(writer, "error sending request", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(writer, "error reading response body", http.StatusInternalServerError)
		return
	}

	if request.Header.Get("Cache-Control") != "no-store" {
		err := setResponseInCache(url, body, response.StatusCode)
		if err != nil {
			fmt.Printf("error setting response in cache for %s. Error: %s\n", url, err.Error())
		}
	} else {
		err := flushCache(url)
		if err != nil {
			fmt.Printf("error removing response in cache for %s. Error: %s\n", url, err.Error())
		}
	}

	writer.WriteHeader(response.StatusCode)
	_, err = writer.Write(body)
	if err != nil {
		http.Error(writer, "error writing response", http.StatusInternalServerError)
		return
	}
}

func getCachedResponse(url string) (ResponseCache, error) {
	redisClient := config.GetRedisClient()

	if redisClient.Ping(context.Background()).Val() == "" {
		return ResponseCache{}, errors.New("could not connect to redis")
	}

	cache, err := redisClient.Get(context.Background(), url).Result()
	if err != nil && err.Error() != "redis: nil" {
		return ResponseCache{}, err
	}

	if cache == "" {
		return ResponseCache{}, nil
	}

	var responseCached ResponseCache

	err = json.Unmarshal([]byte(cache), &responseCached)
	if err != nil && err.Error() != "redis: nil" {
		return ResponseCache{}, err
	}

	return responseCached, nil
}

func setResponseInCache(url string, body []byte, statusCode int) error {
	redisClient := config.GetRedisClient()

	if redisClient.Ping(context.Background()).Val() == "" {
		return errors.New("could not connect to redis")
	}

	cache := ResponseCache{
		StatusCode: statusCode,
		Body:       body,
		Present:    true,
	}

	cacheInJson, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	redisClient.Set(context.Background(), url, cacheInJson, config.GetRedisSecondsToExpireCache())

	return nil
}

func flushCache(url string) error {
	redisClient := config.GetRedisClient()

	if redisClient.Ping(context.Background()).Val() == "" {
		return errors.New("could not connect to redis")
	}

	redisClient.Del(context.Background(), url)

	return nil
}
