package api

import (
	"sync/atomic"

	"github.com/6l20/zama-test/server/usecases"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Router register necessary routes and returns an instance of a router.
func Router(useCases usecases.ServerUseCases) *mux.Router {
	isReady := &atomic.Value{}
	isReady.Store(false)
	
	healthHandler := NewHealthHandler()

	r := mux.NewRouter()
	r.HandleFunc("/upload", useCases.HandleFileUpload()).Methods("POST")
	r.HandleFunc("/download/{filename}", useCases.HandleFileRequest()).Methods("GET")
	r.HandleFunc("/proof/{filenum}", useCases.HandleProofRequest()).Methods("GET")
	r.HandleFunc("/healthz", healthHandler.CheckHealth)
	r.HandleFunc("/readyz", healthHandler.CheckHealth)
	r.Handle("/metrics", promhttp.Handler())
	return r
}