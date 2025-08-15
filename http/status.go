package http

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusMethodNotAllowed    StatusCode = 405
	StatusInternalServerError StatusCode = 500
)

const (
	okMessage                  = "OK"
	badRequestMessage          = "Bad Request"
	notFoundMessage            = "Not Found"
	methodNotAllowedMessage    = "Method Not Allowed"
	internalServerErrorMessage = "Internal Server Error"
)

func (s StatusCode) String() string {
	switch s {
	case StatusOk:
		return okMessage
	case StatusBadRequest:
		return badRequestMessage
	case StatusNotFound:
		return notFoundMessage
	case StatusMethodNotAllowed:
		return methodNotAllowedMessage
	case StatusInternalServerError:
		return internalServerErrorMessage
	default:
		return ""
	}
}
