package http

import (
	"errors"
	"fmt"
)

// ErrMalformedBody is a standard error returned when the server cannot
// parse the request body, for instance, due to chunked encoding errors
// or a mismatched Content-Length.
var ErrMalformedBody = errors.New("body couldn't be parsed")

// HTTPError is a convenience function that sends a standard HTTP error response
// to the client using the status code's default text as the title.
func HTTPError(res Response, status StatusCode) {
	HTTPErrorWithMessage(res, status, status.String(), "")
}

// HTTPErrorWithMessage sends a formatted HTML error page to the client.
// It sets the response status code and includes a title and an optional message.
func HTTPErrorWithMessage(res Response, status StatusCode, title string, message string) {
	res.SetStatus(status)
	res.SetHeader("Content-Type", "text/html")

	messageElm := ""
	if message != "" {
		messageElm = fmt.Sprintf("<p>%s</p>", message)
	}
	res.Write([]byte(fmt.Sprintf("<h1>%s</h1>%s", title, messageElm)))
}
