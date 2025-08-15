package http

import "fmt"

type Version string

const (
	Version10 Version = "HTTP/1.0"
	Version11 Version = "HTTP/1.1"
)

func IsVersionValid(v string) (Version, error) {
	switch v {
	case "HTTP/1.0":
		return Version10, nil
	case "HTTP/1.1":
		return Version11, nil
	default:
		return "", fmt.Errorf("unsupported version %q", v)
	}
}

func (v Version) String() string {
	return string(v)
}
