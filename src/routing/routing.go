package routing

import (
	"errors"
	"go-http-server/src/http"
	"log"
	"regexp"
	"strings"
)

var routingTable = map[string]http.RouteHandlerFunction{
	"^/$":             func(req string) http.HttpResponse { return http.HttpResponse{StatusCode: 200, Body: ""} },
	"^/[a-zA-Z0-9]*$": func(req string) http.HttpResponse { return http.HttpResponse{StatusCode: 200, Body: ""} },
	"^/echo/[a-zA-Z0-9]+$": func(req string) http.HttpResponse {
		// get echo data,
		pattern := "/[a-zA-Z0-9]*$"
		re, err := regexp.Compile(pattern)
		if err != nil {
			return http.HttpResponse{StatusCode: 400, Body: ""}
		}

		matches := re.FindStringSubmatch(req)
		if len(matches) > 0 {
			return http.HttpResponse{StatusCode: 200, Body: strings.ReplaceAll(matches[0], "/", "")}
		}
		return http.HttpResponse{StatusCode: 400, Body: ""}

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
