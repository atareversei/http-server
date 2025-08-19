package http

import (
	"fmt"
	"github.com/atareversei/http-server/internal/qzip"
	"io"
)

type Encoder string

const (
	PLAIN Encoder = "plain"
	QZIP  Encoder = "qzip"
)

func IsEncodingValid(e string) (Encoder, error) {
	switch e {
	case "plain":
		return PLAIN, nil
	case "qzip":
		return QZIP, nil
	}

	return "", fmt.Errorf("unsupported encoding %q", e)
}

func (e Encoder) String() string {
	return string(e)
}

func (e Encoder) NewEncoder(w io.Writer) io.Writer {
	switch e {
	case PLAIN:
		return w
	case QZIP:
		return qzip.NewWriter(w)
	}

	return w
}
