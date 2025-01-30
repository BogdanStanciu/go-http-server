package http

import (
	"bufio"
	"errors"
	"fmt"
	"strings"
)

type HttpResponse struct {
	StatusCode uint16
	Body       string
	Headers    map[string]string
}

type HttpRequest struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    string
	Query   string // Everything after ? or mapping to a map ?
}

func NewHttpResponse() *HttpResponse {
	f := HttpResponse{StatusCode: 0, Body: "", Headers: make(map[string]string)}
	return &f
}
func NewHttpRequest() *HttpRequest {
	f := HttpRequest{Method: "GET", Url: "", Headers: make(map[string]string), Body: "", Query: ""}
	return &f
}

/*
This function will take a raw request and parse it
following the HTTP 1.1 standard, if the request
will not respect the standard will return an error
*/
func HttpParser(rawRequest string) (*HttpRequest, error) {

	if len(rawRequest) == 0 {
		return NewHttpRequest(), errors.New("empty http request")
	}

	request := NewHttpRequest()

	request.Headers = make(map[string]string)

	scanner := bufio.NewScanner(strings.NewReader(rawRequest))
	scanner.Scan()

	line := strings.Split(scanner.Text(), " ")

	request.Method = line[0]
	request.Url = line[1]

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		header := strings.Split(line, ":")
		request.Headers[header[0]] = header[1]
	}

	for scanner.Scan() {
		request.Body += scanner.Text()
	}

	return request, nil

}

/*
This function converts an HTTP response struct
into a properly formatted HTTP response string.
It follows the standard HTTP/1.1 format.
*/
func HttpConverter(request *HttpResponse) (string, error) {
	response := fmt.Sprintf("HTTP/1.1 %d OK \n", request.StatusCode)

	for key, value := range request.Headers {
		response = response + fmt.Sprintf("%s: %s\n", key, value)
	}

	response = response + fmt.Sprintf("\n%s", request.Body)
	return response, nil
}
