package api

import (
	"net/http"
)

type HealthHandler interface {
	CheckHealth(w http.ResponseWriter, r *http.Request)
}

type healthHandler struct {
}

func NewHealthHandler() HealthHandler {
	return &healthHandler{}
}

// swagger:route GET /healthz healthCheck
//
// Health check endpoint
//
//	Produces:
//	- text/plain
//
//	Schemes: http, https
//
//	Responses:
//	  200: OK
//	  503: ServiceUnavailable
func (h *healthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	h.healthz(w, r)
}

// healthz is a liveness probe.
func (h *healthHandler) healthz(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
