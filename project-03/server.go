package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"net"
)

func main() {
	cli.MadeInBasliqLabs()
	portFlag := flag.Int("port", 8080, "port number to spawn the server process")
	flag.Parse()
	port := *portFlag
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		message, err := r.ReadBytes('\n')
		if err != nil {
			cli.Info(fmt.Sprintf("client (%s) closed the connection", conn.RemoteAddr()))
			break
		}
		go cli.Info(fmt.Sprintf("recieved message: %s", message[:len(message)-1]))
		conn.Write(message)
	}
}
