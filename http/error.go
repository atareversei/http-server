package http

import (
	"errors"
	"fmt"
)

var ErrMalformedBody = errors.New("body couldn't be parsed")

func HTTPError(res Response, status StatusCode) {
	HTTPErrorWithMessage(res, status, status.String(), "")
}

func HTTPErrorWithMessage(res Response, status StatusCode, title string, message string) {
	res.WriteHeader(status)
	res.SetHeader("Content-Type", "text/html")

	messageElm := ""
	if message != "" {
		messageElm = fmt.Sprintf("<p>%s</p>", message)
	}
	res.Write([]byte(fmt.Sprintf("<h1>%s</h1>%s", title, messageElm)))
}
