package main

import (
	"fmt"
	"go-http-server/src/http"
	"go-http-server/src/routing"
	"log"
	"net"
	"os"
	"strings"
)

func execRouteHandler(handler *http.RouteHandlerFunction, requestTarget string) string {
	result := (*handler)(requestTarget)

	var message = ""
	if result.StatusCode == 200 {
		message = "OK"
	} else {
		message = "Bad Request"
	}

	response := http.HttpResponseTemplate
	response = strings.Replace(response, "{body}", result.Body, 1)
	response = strings.Replace(response, "{statusCode}", fmt.Sprintf("%d", result.StatusCode), 1)
	response = strings.Replace(response, "{message}", message, 1)
	// calculate Content-Length
	response = strings.Replace(response, "{ContentLength}", fmt.Sprintf("%d", len(result.Body)), 1)

	return response

}

// A request-line begins with a method token
// followed by a single space (SP)
// the request-target, and another single space (SP)
// and ends with the protocol version
func handleRequestLine(request string) (string, error) {

	// split line by SP
	requestParts := strings.Split(request, " ")
	if requestParts[0] != "GET" {
		return "HTTP/1.1 404 Not Found\r\n\r\n", nil
	}

	handler, err := routing.Route(requestParts[1])
	if err != nil {
		return "HTTP/1.1 404 Not Found\r\n\r\n", nil
	}

	return execRouteHandler(handler, request), nil

}

func handleConnection(con net.Conn) {
	defer con.Close()
	req := make([]byte, 1024)
	_, err := con.Read(req)

	if err != nil {
		con.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	}

	response, err := handleRequestLine(string(req))

	if err != nil {
		con.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	}

	con.Write([]byte(response))

}

func main() {
	log.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	log.Println("Server listing on 4221")

	var con net.Conn

	if err != nil {
		log.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	defer func() {
		con.Close()
		l.Close()
	}()

	for {
		con, err := l.Accept()
		log.Printf("Accept conntion from %s\n", con.LocalAddr().String())
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)

		}
		go handleConnection(con)
	}

}
