package main

import (
	"fmt"
	"go-http-server/src/http"
	"go-http-server/src/routing"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Server struct {
	listener net.Listener
	router   *routing.Router
}

func execRouteHandler(handler *routing.RouteHandlerFunction, requestTarget string) string {
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
func (server Server) handleRequestLine(request string) (string, error) {

	// split line by SP
	requestParts := strings.Split(request, " ")
	if requestParts[0] != "GET" {
		return "HTTP/1.1 404 Not Found\r\n\r\n", nil
	}

	handler, err := server.router.Route(requestParts[1])

	if err != nil {
		return "HTTP/1.1 404 Not Found\r\n\r\n", nil
	}

	return execRouteHandler(handler, request), nil
}

func (server Server) handleConnection(con net.Conn) {
	defer con.Close()
	req := make([]byte, 1024)
	_, err := con.Read(req)

	if err != nil {
		con.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	}

	response, err := server.handleRequestLine(string(req))

	if err != nil {
		con.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	}

	con.Write([]byte(response))
}

// Init a new server
func (server *Server) init(port uint16) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	log.Printf("[Server] Server listing on %s\n", fmt.Sprintf("%d", port))

	if err != nil {
		os.Exit(1)
	}
	server.listener = l
	server.router = routing.New()
}

func main() {
	args := os.Args
	if len(args[1:]) < 1 {
		log.Fatal("Please provide a port")
	}
	port, err := strconv.ParseInt(args[1], 10, 16)

	if err != nil {
		log.Fatal(err)
	}

	log.Printf("[Server] Logs from your program will appear here!\n")

	var server Server
	server.init(uint16(port))

	defer func() {
		server.listener.Close()
	}()

	for {
		con, err := server.listener.Accept()

		if err != nil {
			log.Fatalln("[Server] Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go server.handleConnection(con)
	}
}
