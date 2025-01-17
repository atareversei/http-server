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
	"strings"
	"time"
)

func main() {
	fmt.Printf(
		"made in %s",
		colorize.
			New("basliq labs\n").
			Modify(colorize.BrightBlue).
			Commit())

	hostFlag := flag.String("host", "localhost", "Host to dial")
	portFlag := flag.Int("port", 8080, "Port to dial")
	flag.Parse()

	host := *hostFlag
	port := *portFlag

	for {
		fmt.Println("1. measure performance (1 Byte)\n2. measure performance (1 KB)\n3. measure performance (2 KB)\n4. measure performance (4 KB)\n5. measure performance (8 KB)\n6. measure performance (16 KB)\n7. exit")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			cli.Error("invalid input", err)
		}
		input = strings.TrimSuffix(input, "\r\n")
		fmt.Println([]byte(input), []byte("dial"))
		if input == "1" {
			MeasurePerformance(host, port, 1)
		}
		if input == "2" {
			MeasurePerformance(host, port, 1000)
		}
		if input == "3" {
			MeasurePerformance(host, port, 2000)
		}
		if input == "4" {
			MeasurePerformance(host, port, 4000)
		}
		if input == "5" {
			MeasurePerformance(host, port, 8000)
		}
		if input == "6" {
			MeasurePerformance(host, port, 16000)
		}
		if input == "exit" || input == "7" {
			cli.Info("exiting the program")
			os.Exit(1)
		}
	}
}

func MeasurePerformance(host string, port int, n int) {
	times := make([]time.Duration, 0)
	cli.Info(fmt.Sprintf("dialing server at: %s:%d", host, port))
	cli.Info(fmt.Sprintf("sending 1000 messages each %d Byte", n))
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		cli.Error("couldn't connect to the server", err)
		return
	}
	defer conn.Close()
	for i := 0; i < 1000; i++ {
		buf := make([]byte, n+1)
		_, err := rand.Read(buf)
		for i, b := range buf {
			if b == '\n' {
				buf[i] = 'a'
			}
		}
		buf[n] = '\n'
		before := time.Now()
		_, err = conn.Write(buf)
		if err != nil {
			cli.Error("couldn't send the data to server, skipping this iteration", err)
			continue
		}
		resReader := bufio.NewReader(conn)
		_, err = resReader.ReadString('\n')
		if err != nil {
			cli.Error("couldn't read the response", err)
		}
		after := time.Now()
		diff := after.Sub(before)
		times = append(times, diff)
	}
	var sum int64
	for _, t := range times {
		sum += t.Microseconds()
	}
	avg := sum / 1000
	tp := (int64(n) * 8) / (sum / 100000)
	cli.Info(fmt.Sprintf("Measuring finished.\n\t\t\tAverage Time (microseconds): %d\n\t\t\tThroughput (bits/seconds): %d", avg, tp))
}
