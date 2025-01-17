package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/atareversei/network-course-projects/pkg/cli"
	"net"
	"os"
	"strconv"
	"strings"
)

var clientKeys = make(map[int]int)

func main() {
	cli.MadeInBasliqLabs()
	hostFlag := flag.String("host", "localhost", "Host to dial")
	portFlag := flag.Int("port", 8080, "Port to dial")
	flag.Parse()
	host := *hostFlag
	port := *portFlag
	data, err := os.ReadFile("./key.txt")
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
		clientKeys[index] = value
	}
	for {
		fmt.Println("1. dial\n2. exit")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			cli.Error("invalid input", err)
		}
		input = strings.TrimSuffix(input, "\r\n")
		fmt.Println([]byte(input), []byte("dial"))
		if input == "dial" || input == "1" {
			Dial(host, port)
		}
		if input == "exit" || input == "2" {
			cli.Info("exiting the program")
			os.Exit(1)
		}
	}
}

func Dial(host string, port int) {
	cli.Info(fmt.Sprintf("dialing server at: %s:%d", host, port))
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		cli.Error("couldn't connect to the server", err)
		return
	}
	defer conn.Close()
	for {
		fmt.Println("Write a message to be delivered to server or close the connection with `close`:")
		r := bufio.NewReader(os.Stdin)
		input, err := r.ReadString('\n')
		if err != nil {
			cli.Error("couldn't read from Stdin, skipping this iteration", err)
			continue
		}
		input = strings.TrimSuffix(input, "\r\n")
		if input == "close" {
			cli.Info("closing the connection")
			break
		}
		_, err = conn.Write([]byte(input + "\n"))
		if err != nil {
			cli.Error("couldn't send the data to server, skipping this iteration", err)
			continue
		}
		resReader := bufio.NewReader(conn)
		res, err := resReader.ReadString('\n')
		if err != nil {
			cli.Error("couldn't read the response", err)
		}
		res = strings.TrimSuffix(res, "\n")
		dashIndex := strings.Index(res, "-")
		index, err := strconv.Atoi(res[0:dashIndex])
		if err != nil {
			cli.Error(fmt.Sprintf("malformed index: %s", res[0:dashIndex]), err)
			continue
		}
		key := clientKeys[index]
		var builder strings.Builder
		for _, ch := range res[dashIndex+1:] {
			builder.WriteRune(ch - int32(key))
		}
		cli.Success(fmt.Sprintf("encryption was successful:\n\t\t\tresponse: %s\n\t\t\tdecrypted message: %s\n\t\t\tindex used: %d\n\t\t\tkey used: %d", res, builder.String(), index, key))
	}
}
