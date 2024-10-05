package routing

import (
	"errors"
	"go-http-server/src/http"
	"log"
	"os"
	"regexp"
	"strings"
)

var routingTable = map[string]http.RouteHandlerFunction{
	"^/$":             func(req string) http.HttpResponse { return http.HttpResponse{StatusCode: 200, Body: ""} },
	"^/[a-zA-Z0-9]*$": func(req string) http.HttpResponse { return http.HttpResponse{StatusCode: 200, Body: ""} },
	"^/echo/[a-zA-Z0-9]+$": func(req string) http.HttpResponse {

		pattern := "/[a-zA-Z0-9]+$"
		re, err := regexp.Compile(pattern)

		if err != nil {
			return http.HttpResponse{StatusCode: 400, Body: ""}
		}

		matches := re.FindStringSubmatch(strings.Split(req, " ")[1])

		if len(matches) > 0 {
			return http.HttpResponse{StatusCode: 200, Body: strings.ReplaceAll(matches[0], "/", "")}
		}
		return http.HttpResponse{StatusCode: 400, Body: ""}
	},
	"^/file/[a-zA-Z0-9]+$": func(req string) http.HttpResponse {
		pattern := "/[a-zA-Z0-9]+$"
		re, err := regexp.Compile(pattern)

		if err != nil {
			return http.HttpResponse{StatusCode: 404, Body: ""}
		}

		matches := re.FindStringSubmatch(strings.Split(req, " ")[1])
		if len(matches) == 0 {
			return http.HttpResponse{StatusCode: 404, Body: ""}
		}

		basePath, _ := os.Getwd()
		basePath += "/assets/"
		basePath += matches[0]

		data, err := os.ReadFile(basePath)

		if err != nil {
			log.Printf("Failed reading file %s", err)
			return http.HttpResponse{StatusCode: 404, Body: ""}
		}

		return http.HttpResponse{StatusCode: 200, Body: string(data)}
	},
}

// Handle the routing for the given request path
// to the handler function
func Route(requestTarget string) (*http.RouteHandlerFunction, error) {
	for pattern, value := range routingTable {
		match, _ := regexp.MatchString(pattern, requestTarget)
		if match {
			log.Printf("Route found for %s\n", requestTarget)
			return &value, nil
		}
	}
	log.Printf("No route found for %s\n", requestTarget)
	return nil, errors.New("no route found")
}
