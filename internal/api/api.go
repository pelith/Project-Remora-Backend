package api

import (
	"net/http"

	"remora/internal/user"
)

// NewMux builds the HTTP mux for this service (server initialization and dependency injection).
func NewMux(userSvc user.Service) *http.ServeMux {
	mux := http.NewServeMux()
	AddRoutes(mux, userSvc)

	return mux
}
