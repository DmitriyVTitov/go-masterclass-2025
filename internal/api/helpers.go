package api

import (
	"encoding/json"
	"net/http"

	"ugc/internal/api/middleware"
	"ugc/internal/errs"

	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	RequestID  string `json:"requestId"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

// WriteError sends error message and status code.
func (api *API) WriteError(w http.ResponseWriter, r *http.Request, err error) {
	resp := ErrorResponse{
		RequestID: middleware.GetRequestID(r),
		Message:   err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")

	log.Err(err).Send()
	switch err.(type) {
	case errs.ErrNoData:
		resp.StatusCode = http.StatusBadRequest
	case errs.ErrUnauthorized:
		resp.StatusCode = http.StatusUnauthorized
	case errs.ErrBadRequest:
		resp.StatusCode = http.StatusBadRequest
	default:
		resp.StatusCode = http.StatusInternalServerError
	}

	w.WriteHeader(resp.StatusCode)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Err(err).Send()
	}
}

func (api *API) WritePlain(w http.ResponseWriter, r *http.Request, msg string) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte(msg))
	if err != nil {
		api.WriteError(w, r, err)
	}
}

func (api *API) WriteJSON(w http.ResponseWriter, r *http.Request, response any) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		api.WriteError(w, r, err)
	}
}
