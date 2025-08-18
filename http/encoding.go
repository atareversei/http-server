package http

import (
	"fmt"
	"github.com/atareversei/http-server/internal/qzip"
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

func (e Encoder) Encode(in []byte) []byte {
	switch e {
	case PLAIN:
		return in
	case QZIP:
		return qzip.Compress(in)
	}

	return in
}

func (e Encoder) Decode(in []byte) []byte {
	switch e {
	case PLAIN:
		return in
	case QZIP:
		return qzip.Decompress(in)
	}

	return in
}
