package server

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/6l20/zama-test/common/log"

	"github.com/6l20/zama-test/server/config"
)

//go:generate mockgen -source=server.go -destination=mocks/server.go -package=mocks

type IServer interface {
	HandleFileUpload() http.HandlerFunc
	HandleFileRequest() http.HandlerFunc
}

type Server struct {
	Logger log.Logger
	Config config.Config
}

func NewServer(logger log.Logger, config config.Config) *Server {

	
	return &Server{
		Logger: logger,
		Config: config,
	}
}


func (s *Server) HandleFileUpload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("HandleFileUpload")
		const maxUploadSize = 10 * 1024 * 1024 // 10 MB
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, fmt.Sprintf("The uploaded file is too big: %s. Maximum size allowed is %d bytes", err.Error(), maxUploadSize), http.StatusBadRequest)
		return
	}

	// Parse the multipart form, receiving the file.
	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Invalid request. The file field must be named 'uploadFile'.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern.
	tempFile, err := os.Create("/data/" + header.Filename)
	if err != nil {
		http.Error(w, "Failed to create a file", http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	// Read the file's bytes and write them to the temporary file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Invalid file content", http.StatusBadRequest)
		return
	}
	tempFile.Write(fileBytes)

	// Return the name of the file back to the client
	w.Write([]byte(fmt.Sprintf("Successfully uploaded file: %s", header.Filename)))
	}
}

func (s *Server) HandleFileRequest() http.HandlerFunc {

	return nil
}

func (s *Server) SaveFile() error {
	return nil
}
