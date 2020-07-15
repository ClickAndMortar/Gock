package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var stdLog = logrus.New()
var errLog = logrus.New()

func main() {
	stdLog.Out = os.Stdout
	debug := getEnvDefault("GOCK_DEBUG", "0")
	if debug == "1" {
		stdLog.SetLevel(logrus.DebugLevel)
	}
	errLog.Out = os.Stderr

	port := getEnvDefault("GOCK_PORT", "8000")

	mode := getEnvDefault("GOCK_MODE", "default")
	if mode == "proxy" {
		http.HandleFunc("/", ProxyRequestHandler)
	} else {
		http.HandleFunc("/", DefaultRequestHandler)
	}

	rand.Seed(time.Now().UnixNano())

	stdLog.Infof("Listening on http://0.0.0.0:%s in %s mode", port, mode)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}

func ProxyRequestHandler(writer http.ResponseWriter, request *http.Request) {
	percent, _ := strconv.Atoi(getEnvDefault("GOCK_PROXY_PERCENT", "100"))
	returnCode := 0

	random := rand.Intn(100)
	if random < percent {
		stdLog.Debugf("Fault injection: random number %d < percentage %d", random, percent)
		wait := getEnvDefault("GOCK_PROXY_WAIT", "0")
		waitTime, _ := strconv.Atoi(wait)
		duration, _ := time.ParseDuration(fmt.Sprintf("%ds", waitTime))
		time.Sleep(duration)

		code := getEnvDefault("GOCK_PROXY_CODE", "0")
		if code != "" {
			returnCode, _ = strconv.Atoi(code)
		}
	}

	if request.URL.Scheme != "http" && request.URL.Scheme != "https" {
		request.URL.Scheme = "http"
	}

	stdLog.Printf("%s %s", request.Method, strings.ReplaceAll(request.RequestURI, "%", "%%"))

	client := &http.Client{}

	request.URL.Host = fmt.Sprintf("%s:%s", getEnvDefault("GOCK_PROXY_HOST", ""), getEnvDefault("GOCK_PROXY_PORT", "80"))
	request.RequestURI = ""

	delHopHeaders(request.Header)

	if clientIP, _, err := net.SplitHostPort(request.RemoteAddr); err == nil {
		appendHostToXForwardHeader(request.Header, clientIP)
	}

	resp, err := client.Do(request)
	if err != nil {
		http.Error(writer, "Bad Gateway", http.StatusBadGateway)
		errLog.Errorf("Error proxying backend: %s", strings.ReplaceAll(err.Error(), "%", "%%"))
		return
	}
	defer resp.Body.Close()

	delHopHeaders(resp.Header)

	copyHeader(writer.Header(), resp.Header)
	if returnCode != 0 {
		writer.WriteHeader(returnCode)
	} else {
		writer.WriteHeader(resp.StatusCode)
	}
	io.Copy(writer, resp.Body)
}

func DefaultRequestHandler(writer http.ResponseWriter, request *http.Request) {

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

	stdLog.Printf("%s %s", request.Method, request.RequestURI)

	writer.WriteHeader(returnCode)
	writer.Write([]byte(""))
}

func getEnvDefault(name string, defaultValue string) string {
	variable := os.Getenv(name)
	if variable == "" && defaultValue != "" {
		stdLog.Infof("Environment variable %s not set or empty, using default %s", name, defaultValue)
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

var hopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te", // canonicalized version of "TE"
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func delHopHeaders(header http.Header) {
	for _, h := range hopHeaders {
		header.Del(h)
	}
}

func appendHostToXForwardHeader(header http.Header, host string) {
	// If we aren't the first proxy retain prior
	// X-Forwarded-For information as a comma+space
	// separated list and fold multiple headers into one.
	if prior, ok := header["X-Forwarded-For"]; ok {
		host = strings.Join(prior, ", ") + ", " + host
	}
	header.Set("X-Forwarded-For", host)
}
