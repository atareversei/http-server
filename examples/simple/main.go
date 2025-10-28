package main

import (
	"flag"

	"github.com/atareversei/http-server/examples/simple/handler"
	"github.com/atareversei/http-server/examples/simple/repository"
	"github.com/atareversei/http-server/http"
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

	s := http.New()
	router := s.NewRouter()
	router.Get("/contacts",  app.handler.Contact)
	s.Start(port)
}
