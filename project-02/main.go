package main

import (
	"flag"
	"github.com/atareversei/network-course-projects/pkg/http/server"
	"github.com/atareversei/network-course-projects/project-02/handler"
	"github.com/atareversei/network-course-projects/project-02/repository"
)

type App struct {
	handler handler.Handler
}

func main() {
	portFlag := flag.Int("port", 8080, "Port to serve")
	flag.Parse()
	port := *portFlag

	repo := repository.New()
	hndlr := handler.New(&repo)
	app := App{handler: hndlr}

	s := server.New(port)
	s.Get("/contacts", app.handler.Contact)
	s.Start()
}
