package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", RequestHandler)

	port := getEnvDefault("GOCK_PORT", "8000")

	log.Printf("Listening on http://0.0.0.0:%s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func RequestHandler(writer http.ResponseWriter, request *http.Request) {

	wait := request.URL.Query().Get("wait")
	if wait != "" {
		duration, _ := time.ParseDuration(fmt.Sprintf("%ss", wait))
		time.Sleep(duration)
	}

	code := request.URL.Query().Get("code")
	returnCode := 204
	if code != "" {
		returnCode, _ = strconv.Atoi(code)
	}

	random := request.URL.Query().Get("random")
	randoms := strings.Split(random, ",")
	if contains(randoms, "code") {

	}

	writer.WriteHeader(returnCode)
	writer.Write([]byte(""))
}

func getEnvDefault(name string, defaultValue string) string {
	variable := os.Getenv(name)
	if variable == "" && defaultValue != "" {
		log.Printf("Environment variable %s not set or empty, using default %s", name, defaultValue)
		variable = defaultValue
	}

	return variable
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
