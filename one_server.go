package main

import (
	"flag"
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"github.com/atareversei/network-course-projects/pkg/colorize"
	"net"
	"os"
	"strconv"
	"strings"
)

var keys = make(map[int]int)

func main() {
	fmt.Printf(
		"made in %s",
		colorize.
			New("basliq labs\n").
			Modify(colorize.BrightBlue).
			Commit())

	port := flag.Int("port", 8080, "port number to spawn the server process")
	flag.Parse()

	data, err := os.ReadFile("./one_key.txt")
	if err != nil {
		cli.Error("server could not find the keys", err)
		os.Exit(1)
	}

	for _, key := range strings.Split(strings.ReplaceAll(string(data), "\r", ""), "\n") {
		index, err := strconv.Atoi(strings.Trim(strings.Split(key, ":")[0], " "))
		if err != nil {
			cli.Error("cannot read key index", err)
			continue
		}

		value, err := strconv.Atoi(strings.Trim(strings.Split(key, ":")[1], " "))
		if err != nil {
			cli.Error("cannot read key value", err)
			continue
		}
		keys[index] = value
	}

	cli.Success(fmt.Sprintf("tcp server started at :%d", *port))
	skt, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		cli.Error("server could not be started", err)
	}

	for {
		conn, err := skt.Accept()
		if err != nil {
			cli.Error("connection resulted in an error", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		cli.Error("could not read the content", err)
		return
	}
	r := string(buf)
	cli.Info(fmt.Sprintf("recieved bytes: %v", r))
}
