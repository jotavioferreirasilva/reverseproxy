package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type numbersToSum struct {
	Numbers string `json:"numbers"`
}

func Ping(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		_, err := writer.Write([]byte("pong"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func Sum(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		numbers, present := request.URL.Query()["number"]
		if !present {
			http.Error(writer, "numbers must be informed", http.StatusBadRequest)
			return
		}
		sum, err := sumNumbers(numbers)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = writer.Write([]byte(fmt.Sprintf("Sum is equal to:%d", sum)))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		body, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "error reading body", http.StatusBadRequest)
			return
		}
		var numbersToSum numbersToSum

		err = json.Unmarshal(body, &numbersToSum)
		if err != nil {
			http.Error(writer, "error reading body", http.StatusBadRequest)
			return
		}

		sum, err := sumNumbers(strings.Split(numbersToSum.Numbers, ","))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = writer.Write([]byte(fmt.Sprintf("Sum is equal to:%d", sum)))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func sumNumbers(numbers []string) (int, error) {
	sum := 0
	for _, number := range numbers {
		convertedNumber, err := strconv.Atoi(number)
		if err != nil {
			return 0, errors.New("parameter number is is invalid format")
		}
		sum += convertedNumber
	}
	return sum, nil
}
