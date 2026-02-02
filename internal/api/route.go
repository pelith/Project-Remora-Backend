package api

import (
	"net/http"

	"remora/internal/user"
	userapi "remora/internal/user/api"
)

// AddRoutes registers API routes on the provided mux (central routing).
func AddRoutes(mux *http.ServeMux, userSvc user.Service) {
	userapi.AddRoutes(mux, userSvc)
}
