package main

import "github.com/atareversei/http-server/http"

func main() {
	server := http.New(9012, nil)
	server.Start()
}
