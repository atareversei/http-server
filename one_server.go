package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"github.com/atareversei/network-course-projects/pkg/colorize"
	"math/rand"
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

	portFlag := flag.Int("port", 8080, "port number to spawn the server process")
	flag.Parse()

	port := *portFlag

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
		message, err := r.ReadString('\n')
		if err != nil {
			cli.Info(fmt.Sprintf("client (%s) closed the connection", conn.RemoteAddr()))
			break
		}
		message = strings.TrimSuffix(message, "\n")
		cli.Info(fmt.Sprintf("recieved message: %s", message))
		ucm := strings.ToUpper(message)
		index := rand.Intn(10)
		key := keys[index]

		var builder strings.Builder
		for _, ch := range ucm {
			builder.WriteRune(ch + int32(key))
		}
		response := fmt.Sprintf("%d-%s\n", index, builder.String())
		cli.Success(fmt.Sprintf("encryption was successful:\n\t\t\tmessage: %s\n\t\t\tencrypted message: %s\n\t\t\tindex used: %d\n\t\t\tkey used: %d", message, builder.String(), index, key))
		conn.Write([]byte(response))
	}
}
