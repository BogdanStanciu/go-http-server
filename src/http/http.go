package http

type HttpResponse struct {
	StatusCode uint16
	Body       string
}

var HttpResponseTemplate string = "HTTP/1.1 {statusCode} {message}\r\nContent-Type: text/plain\r\nContent-Length: {ContentLength}\r\n\r\n{body}"
