package http

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusNoContent           StatusCode = 204
	StatusBadRequest          StatusCode = 400
	StatusNotFound            StatusCode = 404
	StatusMethodNotAllowed    StatusCode = 405
	StatusInternalServerError StatusCode = 500
	StatusServiceUnavailable  StatusCode = 503
)

const (
	okMessage                  = "OK"
	noContentMessage           = "No Content"
	badRequestMessage          = "Bad Request"
	notFoundMessage            = "Not Found"
	methodNotAllowedMessage    = "Method Not Allowed"
	internalServerErrorMessage = "Internal Server Error"
	serviceUnavailableMessage  = "Service Unavailable"
)

func (s StatusCode) String() string {
	switch s {
	case StatusOk:
		return okMessage
	case StatusNoContent:
		return noContentMessage
	case StatusBadRequest:
		return badRequestMessage
	case StatusNotFound:
		return notFoundMessage
	case StatusMethodNotAllowed:
		return methodNotAllowedMessage
	case StatusInternalServerError:
		return internalServerErrorMessage
	case StatusServiceUnavailable:
		return serviceUnavailableMessage
	default:
		return ""
	}
}
