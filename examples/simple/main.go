package main

import "github.com/atareversei/http-server/http"

func main() {
	server := http.Server{Port: 9012}
	server.Start()
}
