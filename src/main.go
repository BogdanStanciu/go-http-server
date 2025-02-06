package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"go-http-server/src/http"
	"go-http-server/src/routing"
	"log"
	"net"
	"os"
	"strconv"
)

type Server struct {
	listener net.Listener
	router   *routing.Router
}

/*
Compress data with gzip encoding
*/
func compress(data string) (string, error) {
	var buf bytes.Buffer
	compressor := gzip.NewWriter(&buf)

	_, err := compressor.Write([]byte(data))

	if err != nil {
		return "", err
	}

	if err := compressor.Close(); err != nil {
		return "", err
	}

	return buf.String(), nil
}

/*
For this first version the server will handle only GET request.
This function check if the incoming request is a GET and if the route exits
then exec the associated handler
*/
func (server Server) execRouteHandler(request *http.HttpRequest) (string, error) {

	if request.Method != "GET" {
		return "HTTP/1.1 404 Not Found\r\n\r\n", nil
	}

	handler, err := server.router.Route(request.Url)

	if err != nil {
		return "HTTP/1.1 404 Not Found\r\n\r\n", nil
	}

	result := (*handler)(request)

	if request.Headers["Accept-Encoding"] != "" {
		result.Body, _ = compress(result.Body)
		result.Headers["Content-Length"] = fmt.Sprintf("%d", len(result.Body))
		result.Headers["Content-Encoding"] = "gzip"
	}

	res, err := http.HttpConverter(result)

	if err != nil {
		return "HTTP/1.1 400 Bad Request\r\n\r\n\r\n", nil
	}

	return res, nil
}

/*
Handle the init of the connection and ready the incoming bytes.
After read all the incoming request will try to parse as a HTTP 1.1 request
then will handle the request, else will reponde with 400
*/
func (server Server) handleConnection(con net.Conn) {
	defer con.Close()

	// handle requests of 1mb
	req := make([]byte, 1024*1024)
	n, err := con.Read(req)

	if n == 1024*1024 {
		con.Write([]byte("HTTP/1.1 413 Content Too Large\r\n\r\n"))
		return
	}

	if err != nil || n == 0 {
		con.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	httpRequest, err := http.HttpParser(string(req))

	if err != nil {
		con.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	response, err := server.execRouteHandler(httpRequest)

	if err != nil {
		con.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	con.Write([]byte(response))
}

/*
Init a new server listening to a given port
*/
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
