package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"net"
	"os"
	"strconv"
	"strings"
)

type Contact struct {
	code    string
	name    string
	phone   string
	address string
	email   string
}

type HTTPRequest struct {
	method   string
	resource string
	version  string
	headers  map[string]string
	body     []byte
	conn     net.Conn
}

var MalformedBodyError = errors.New("body couldn't be parsed")

func main() {
	portFlag := flag.Int("port", 8080, "Port to serve")
	flag.Parse()
	port := *portFlag

	contacts := make([]Contact, 0)
	data, err := os.ReadFile("./two_data.txt")
	if err != nil {
		cli.Error("server could not find the data file", err)
		os.Exit(1)
	}

	for _, record := range strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n") {
		contactInfo := strings.Split(record, ",")
		var contact Contact
		contact.code = contactInfo[0]
		contact.name = contactInfo[1]
		contact.phone = contactInfo[2]
		contact.address = contactInfo[3]
		contact.email = contactInfo[4]

		contacts = append(contacts, contact)
	}

	cli.Success(fmt.Sprintf("tcp server started at :%d", port))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		cli.Error("server could not be started", err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			cli.Error("connection resulted in an error", err)
			continue
		}
		go parseHTTP(conn)
	}
}

func parseHTTP(conn net.Conn) (HTTPRequest, error) {
	defer conn.Close()
	var httpRequest HTTPRequest
	httpRequest.conn = conn
	r := bufio.NewReader(conn)
	request, err := r.ReadString('\n')
	if err != nil {
		cli.Error("failed to read the request", err)
		return HTTPRequest{}, err
	}
	requestParts := strings.Split(request, " ")
	httpRequest.method = strings.TrimSpace(requestParts[0])
	httpRequest.resource = strings.TrimSpace(requestParts[1])
	httpRequest.version = strings.TrimSpace(requestParts[2])

	httpRequest.headers = make(map[string]string)
	for {
		header, err := r.ReadString('\n')
		header = strings.TrimSpace(header)
		fmt.Println(header)
		if err != nil {
			cli.Error(fmt.Sprintf("couldn't read header, skipping header (%s)", header), err)
			continue
		}
		// empty line which indicates the header section is over
		if header == "" {
			break
		}
		headerParts := strings.SplitN(header, ":", 2)
		if len(headerParts) != 2 {
			cli.Error(fmt.Sprintf("malformed header: %s", header), fmt.Errorf("not abiding by key:value format"))
			continue
		}
		httpRequest.headers[headerParts[0]] = headerParts[1]
	}
	contentLengthValue, ok := httpRequest.headers["Content-Length"]
	if !ok {
		cli.Error("couldn't find `Content-Length` in the headers", fmt.Errorf("headers map returned `!ok` for `Content-Length`"))
		return httpRequest, MalformedBodyError
	}
	contentLength, err := strconv.Atoi(contentLengthValue)
	if err != nil {
		cli.Error("couldn't parse `Content-Length` out of headers", err)
		return httpRequest, MalformedBodyError
	}
	if contentLength == 0 {
		return httpRequest, nil
	}
	body := make([]byte, contentLength)
	_, err = r.Read(body)
	httpRequest.body = body
	fmt.Println("resource: ", httpRequest.resource)
	return httpRequest, nil
}
