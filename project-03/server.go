package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"github.com/atareversei/network-course-projects/pkg/colorize"
	"net"
)

func main() {
	fmt.Printf(
		"made in %s",
		colorize.
			New("basliq labs\n").
			Modify(colorize.BrightBlue).
			Commit())

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
		cli.Info(fmt.Sprintf("recieved message: %s", message))
		conn.Write(message)
	}
}
