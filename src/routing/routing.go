package routing

import (
	"errors"
	"go-http-server/src/http"
	"log"
	"os"
	"regexp"
	"strings"
)

type RouteHandlerFunction func(request http.HttpRequest) http.HttpResponse

type Router struct {
	tree *Node
}

func New() *Router {
	router := Router{}
	router.tree = InitTree()

	// add routes
	router.tree.AddNode(strings.Split("", "/")[:1], func(req http.HttpRequest) http.HttpResponse { return http.HttpResponse{StatusCode: 200, Body: ""} })
	router.tree.AddNode(strings.Split("/hello", "/")[1:], func(req http.HttpRequest) http.HttpResponse {
		return http.HttpResponse{StatusCode: 200, Body: "Hello, World !"}
	})
	router.tree.AddNode(strings.Split("/echo/*", "/")[1:], func(req http.HttpRequest) http.HttpResponse {
		pattern := "/[a-zA-Z0-9]+$"
		re, err := regexp.Compile(pattern)

		if err != nil {
			return http.HttpResponse{StatusCode: 400, Body: ""}
		}

		matches := re.FindStringSubmatch(req.Url)

		if len(matches) > 0 {
			return http.HttpResponse{StatusCode: 200, Body: strings.ReplaceAll(matches[0], "/", "")}
		}
		return http.HttpResponse{StatusCode: 400, Body: ""}
	})

	router.tree.AddNode(strings.Split("/file/*", "/")[1:], func(req http.HttpRequest) http.HttpResponse {
		pattern := "/[a-zA-Z0-9]+$"
		re, err := regexp.Compile(pattern)

		if err != nil {
			return http.HttpResponse{StatusCode: 404, Body: ""}
		}

		matches := re.FindStringSubmatch(req.Url)

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
	})

	router.tree.PrintTree("")
	return &router

}

func (router *Router) Route(request string) (*RouteHandlerFunction, error) {
	path := strings.Split(request, "/")
	node := router.tree.Search(path[1:])

	if node == nil || !node.IsLeaf {
		return nil, errors.New("no route found")
	}

	return &node.Handler, nil
}
