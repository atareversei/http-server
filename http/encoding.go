package http

import (
	"fmt"
	"github.com/atareversei/http-server/internal/qzip"
	"io"
)

// Encoder represents a content encoding type, like "qzip".
type Encoder string

const (
	PLAIN Encoder = "plain" // PLAIN represents no encoding. The data is sent as-is.
	QZIP  Encoder = "qzip"  // QZIP represents an educational qzip (LZ77) compression.
)

// IsEncodingValid checks if a string corresponds to a supported encoding type.
// It returns the corresponding Encoder and a nil error if valid.
func IsEncodingValid(e string) (Encoder, error) {
	switch e {
	case "plain":
		return PLAIN, nil
	case "qzip":
		return QZIP, nil
	}

	return "", fmt.Errorf("unsupported encoding %q", e)
}

// String returns the string representation of the Encoder.
func (e Encoder) String() string {
	return string(e)
}

// NewEncoder acts as a factory, returning an io.Writer that wraps the provided
// writer with the specified encoding logic.
func (e Encoder) NewEncoder(w io.Writer) io.Writer {
	switch e {
	case PLAIN:
		return w
	case QZIP:
		return qzip.NewWriter(w)
	}

	return w
}
