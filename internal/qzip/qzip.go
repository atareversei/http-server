package qzip

import "io"

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
	return &Writer{w: w, window: make([]byte, 4096)}
}

func (w *Writer) Write(data []byte) (int, error) {
	w.window = append(w.window, data...)
	processed := 0

	for len(w.window)-processed >= maxMatchLength {
		//distance, length := w.findLongestMatch(processed)
	}

	return 0, nil
}

// findLongestMatch searches backwards in the window for the longest repeating sequence.
func (w *Writer) findLongestMatch(curPos int) (int, int) {
	bestLength := 0
	bestDistance := 0

	lookahead := w.window[curPos:]

	for distance := 1; distance <= curPos; distance++ {
		currentLength := 0

		for i := 0; i < len(lookahead); i++ {
			if w.window[curPos-distance+i] == lookahead[i] {
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

func Compress(raw []byte) []byte {
	return raw
}

func Decompress(enc []byte) []byte {
	return enc
}
