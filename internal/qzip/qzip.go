package qzip

import (
	"bufio"
	"errors"
	"io"
)

const (
	literalMarker byte = 0
	pointerMarker byte = 1

	maxMatchLength int = 3
)

// Writer is a streaming writer that performs a simplified LZ77 compression.
type Writer struct {
	w      io.Writer
	window []byte
}

// NewWriter creates a new qzip Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w, window: make([]byte, 0, 4096)}
}

func (qw *Writer) Write(data []byte) (int, error) {
	qw.window = append(qw.window, data...)
	processed := 0

	for len(qw.window)-processed >= maxMatchLength {
		distance, length := qw.findLongestMatch(processed)
		if length > 3 {
			_, err := qw.w.Write([]byte{pointerMarker, byte(distance), byte(length)})
			if err != nil {
				return processed, err
			}
			processed += length
		} else {
			_, err := qw.w.Write([]byte{literalMarker, qw.window[processed]})
			if err != nil {
				return processed, err
			}
			processed++
		}
	}

	return 0, nil
}

// findLongestMatch searches backwards in the window for the longest repeating sequence.
func (qw *Writer) findLongestMatch(curPos int) (int, int) {
	bestLength := 0
	bestDistance := 0

	lookahead := qw.window[curPos:]

	for distance := 1; distance <= curPos; distance++ {
		currentLength := 0

		for i := 0; i < len(lookahead); i++ {
			if qw.window[curPos-distance+i] == lookahead[i] {
				currentLength++
			} else {
				break
			}
		}

		if currentLength > bestLength {
			bestLength = currentLength
			bestDistance = distance
		}
	}

	return bestDistance, bestLength
}

func (qw *Writer) Close() error {
	for _, b := range qw.window {
		_, err := qw.w.Write([]byte{literalMarker, b})
		if err != nil {
			return err
		}
	}
	qw.window = nil
	return nil
}

type Reader struct {
	r      *bufio.Reader
	window []byte
	outbuf []byte
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: bufio.NewReader(r), window: make([]byte, 0, 4096)}
}

func (qr *Reader) Read(data []byte) (int, error) {
	bytesRead := 0

	for bytesRead < len(data) {
		if len(qr.outbuf) > 0 {
			n := copy(data[bytesRead:], qr.outbuf)
			bytesRead += n
			qr.outbuf = qr.outbuf[n:]
			continue
		}
		marker, err := qr.r.ReadByte()
		if err != nil {
			return bytesRead, err
		}
		if marker == literalMarker {
			literal, err := qr.r.ReadByte()
			if err != nil {
				return bytesRead, err
			}
			data[bytesRead] = literal
			bytesRead++
			qr.window = append(qr.window, literal)
		} else if marker == pointerMarker {
			distByte, err := qr.r.ReadByte()
			if err != nil {
				return bytesRead, err
			}
			lenByte, err := qr.r.ReadByte()
			if err != nil {
				return bytesRead, err
			}
			distance := int(distByte)
			length := int(lenByte)

			start := len(qr.window) - distance
			if start < 0 {
				return bytesRead, errors.New("qzip: invalid distance in pointer")
			}

			sequence := qr.window[start : start+length]
			n := copy(data[bytesRead:], sequence)
			bytesRead += n

			if n < len(sequence) {
				qr.outbuf = sequence[n:]
			}
			qr.window = append(qr.window, sequence...)
		} else {
			return bytesRead, errors.New("qzip: invalid marker in data stream")
		}
	}

	return bytesRead, nil
}
