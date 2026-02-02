package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"

	"remora/internal/httpwrap"
	"remora/internal/user"
)

// AddRoutes registers user-related routes on the provided mux.
func AddRoutes(mux *http.ServeMux, svc user.Service) {
	mux.HandleFunc("GET /users/{id}", httpwrap.Handler(getUser(svc)))
}

// UserResponse is the API response for a single user (api-guide: Response types use suffix Response).
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// getUser returns a handler that fetches a user by ID. It returns (*Response, *ErrorResponse) per api-guide Response Construction.
func getUser(svc user.Service) func(*http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
	return func(r *http.Request) (*httpwrap.Response, *httpwrap.ErrorResponse) {
		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			return nil, httpwrap.NewInvalidParamErrorResponse("id")
		}

		u, err := svc.ByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, user.ErrNotFound) {
				return nil, &httpwrap.ErrorResponse{
					StatusCode: http.StatusNotFound,
					ErrorMsg:   "not found",
					Err:        err,
				}
			}

			return nil, &httpwrap.ErrorResponse{
				StatusCode: http.StatusInternalServerError,
				ErrorMsg:   "internal error",
				Err:        err,
			}
		}

		return &httpwrap.Response{
			StatusCode: http.StatusOK,
			Body: &UserResponse{
				ID:        u.ID,
				Address:   u.Address,
				CreatedAt: u.CreatedAt,
				UpdatedAt: u.UpdatedAt,
			},
		}, nil
	}
}
