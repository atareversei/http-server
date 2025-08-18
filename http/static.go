package http

import (
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// FileHandler maps a URL path prefix to a directory on disk for serving static files.
func (s *Server) FileHandler(pattern string, directory string) {
	s.Static[pattern] = directory
}

// handleFileRequest serves files from a static directory mapped by the given prefix.
// It sets proper MIME types and handles 404/500 errors.
func (s *Server) handleFileRequest(prefix string, request Request, response Response) {
	i := strings.Index(request.Path(), prefix)
	filePath := request.Path()[i+len(prefix):]
	fullPath := s.Static[prefix] + filePath
	f, err := os.Open(fullPath)
	defer f.Close()

	if os.IsNotExist(err) {
		response.SetStatus(StatusNotFound)
		response.SetHeader("Content-Type", "text/html")
		response.Write([]byte("<h1>404 Not Found</h1>"))
		return
	} else if err != nil {
		response.SetStatus(StatusInternalServerError)
		response.SetHeader("Content-Type", "text/html")
		response.Write([]byte("<h1>500 Internal Server Error</h1>"))
		return
	}

	fileInfo, err := f.Stat()
	if err != nil {
		HTTPError(response, StatusInternalServerError)
		return
	}
	response.contentLength = int(fileInfo.Size())

	ext := filepath.Ext(fullPath)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	response.SetHeader("Content-Type", mimeType)

	_, err = io.Copy(response.conn, f)
	if err != nil {
		response.conn.Close()
	}
}
