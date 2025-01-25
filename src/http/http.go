package http

import (
	"bufio"
	"errors"
	"strings"
)

type HttpResponse struct {
	StatusCode uint16
	Body       string
}

type HttpRequest struct {
	Method  string
	Url     string
	Headers map[string]string
	Body    string
	Query   string // Everything after ? or mapping to a map ?
}

var HttpResponseTemplate string = "HTTP/1.1 {statusCode} {message}\r\nContent-Type: text/plain\r\nContent-Length: {ContentLength}\r\n\r\n{body}"

/*
This function will take a raw request and parse it
following the HTTP 1.1 standard, if the request
will not respect the standard will return an error
*/
func HttpParser(rawRequest string) (HttpRequest, error) {
	// read line by line
	// first line is request path and uri
	// method SP request-target SP HTTP-version

	if len(rawRequest) == 0 {
		return HttpRequest{}, errors.New("empty http request")
	}

	var request HttpRequest
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
