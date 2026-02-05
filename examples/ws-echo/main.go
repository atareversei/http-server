package main

import (
	"flag"
	"net"

	"github.com/atareversei/http-server/examples/ws-echo/handler"
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

	websocketServer := ws.Server{}

	hndlr := handler.New()
	app := App{handler: hndlr}
	s := http.New()
	router := s.NewRouter()
	router.Get("/ping", app.handler.Ping)

	s.UpgradeHandler = func(conn *net.Conn, req http.Request) {
		websocketServer.Start(
			conn,
			ws.HTTPRequest{
				Path:    req.Path(),
				Method:  req.Method().String(),
				Headers: req.Headers(),
				Params:  req.Params(),
			},
		)
	}

	s.Start(port)
}
