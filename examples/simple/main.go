package main

import "github.com/atareversei/http-server/http"

func main() {

	server := http.New(9012)
	server.NewRouter().Get("/contacts", func(req http.Request, res http.Response) {})
	server.Start()
}
