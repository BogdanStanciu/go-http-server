package routing

import (
	"errors"
	"go-http-server/src/http"
	"log"
	"os"
	"regexp"
	"strings"
)

type RouteHandlerFunction func(request *http.HttpRequest) *http.HttpResponse

type Router struct {
	tree *Node
}

func New() *Router {
	router := Router{}
	router.tree = InitTree()

	// add routes
	router.tree.AddNode(strings.Split("", "/")[:1], func(req *http.HttpRequest) *http.HttpResponse {
		res := http.NewHttpResponse()
		res.StatusCode = 200
		res.Body = ""
		return res
	})

	router.tree.AddNode(strings.Split("/hello", "/")[1:], func(req *http.HttpRequest) *http.HttpResponse {
		res := http.NewHttpResponse()
		res.StatusCode = 200
		res.Body = "Hello World"
		return res
	})

	router.tree.AddNode(strings.Split("/echo/*", "/")[1:], func(req *http.HttpRequest) *http.HttpResponse {
		pattern := "/[a-zA-Z0-9]+$"
		re, err := regexp.Compile(pattern)

		if err != nil {
			res := http.NewHttpResponse()
			res.StatusCode = 400
			res.Body = ""
			return res
		}

		matches := re.FindStringSubmatch(req.Url)

		res := http.NewHttpResponse()

		if len(matches) > 0 {
			res.StatusCode = 200
			res.Body = strings.ReplaceAll(matches[0], "/", "")
			return res

		}
		res.Body = ""
		res.StatusCode = 400
		return res
	})

	router.tree.AddNode(strings.Split("/file/*", "/")[1:], func(req *http.HttpRequest) *http.HttpResponse {
		pattern := "/[a-zA-Z0-9]+$"
		re, err := regexp.Compile(pattern)

		if err != nil {
			res := http.NewHttpResponse()
			res.StatusCode = 404
			res.Body = ""
			return res

		}

		matches := re.FindStringSubmatch(req.Url)

		if len(matches) == 0 {
			res := http.NewHttpResponse()
			res.StatusCode = 404
			res.Body = ""
			return res
		}

		basePath, _ := os.Getwd()
		basePath += "/assets/"
		basePath += matches[0]

		data, err := os.ReadFile(basePath)

		res := http.NewHttpResponse()
		if err != nil {
			log.Printf("Failed reading file %s", err)
			res.StatusCode = 400
			return res
		}

		res.StatusCode = 200
		res.Body = string(data)
		return res
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
