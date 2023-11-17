package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/6l20/zama-test/common/log"
	"github.com/6l20/zama-test/common/merkle"
	"github.com/gorilla/mux"

	"github.com/6l20/zama-test/server/config"
)

//go:generate mockgen -source=server.go -destination=mocks/server.go -package=mocks

type IServer interface {
	HandleFileUpload() http.HandlerFunc
	HandleFileRequest() http.HandlerFunc
	HandleProofRequest() http.HandlerFunc
}

type Server struct {
	Logger log.Logger
	Config config.Config
	merkleManager *merkle.MerkleManager
}

func NewServer(logger log.Logger, config config.Config) *Server {

	
	return &Server{
		Logger: logger,
		Config: config,
		merkleManager: merkle.NewMerkleManager(logger.WithComponent("merkle-server")),
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

	return func (w http.ResponseWriter, r *http.Request) {
		// Extract the file name from the URL query parameters
		fileName := r.URL.Query().Get("filename")
		if fileName == "" {
			http.Error(w, "Filename is required", http.StatusBadRequest)
			return
		}
	
		// Specify the directory where files are stored
		fileDir := "path/to/your/files/"
	
		// Open the file
		file, err := os.Open(fileDir + fileName)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer file.Close()
	
		// Set the correct headers for content-disposition and content type
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")
	
		// Stream the file content to the response
		io.Copy(w, file)
	}
}

func (s *Server) HandleProofRequest() http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		// Extract the file name from the URL query parameters
		s.Logger.Info("HandleProofRequest")
		vars := mux.Vars(r)
		fileNumber := vars["filenum"]
		if fileNumber == "" {
			http.Error(w, "filenum is required", http.StatusBadRequest)
			return
		}

		f, err:= strconv.Atoi(fileNumber) 
		if err != nil {
			http.Error(w, "filenum must be an integer", http.StatusBadRequest)
			return
		}
		
		s.merkleManager.BuildMerkleTreeFromFS("/data")
		proof, err := s.merkleManager.GenerateProof(f)
		if err != nil {
			http.Error(w, "Error generating proof", http.StatusBadRequest)
			return
		}

		if proof == nil {
			http.Error(w, "Proof is nil", http.StatusBadRequest)
			return
		}

		// Send the proof back to the client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(proof)
		
	}
}

func (s *Server) SaveFile() error {
	return nil
}
