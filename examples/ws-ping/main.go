package main

import (
	"flag"
	"net"

	"github.com/atareversei/http-server/examples/ws-ping/handler"
	"github.com/atareversei/http-server/http"
	"github.com/atareversei/http-server/ws"
)

type App struct {
	handler handler.Handler
}

func main() {
	portFlag := flag.Int("port", 8080, "Port to serve")
	flag.Parse()
	port := *portFlag

	ws := ws.Server{}

	hndlr := handler.New()
	app := App{handler: hndlr}
	s := http.New()
	router := s.NewRouter()
	router.Get("/ping", app.handler.Ping)

	s.UpgradeHandler = func(conn *net.Conn, req http.Request) {
		ws.Start(conn)
	}

	s.Start(port)
}
