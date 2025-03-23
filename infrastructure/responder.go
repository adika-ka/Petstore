package infrastructure

import (
	"encoding/json"
	"net/http"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorNotFound(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type JSONResponder struct{}

func NewJSONResponder() Responder {
	return &JSONResponder{}
}

func (r *JSONResponder) OutputJSON(w http.ResponseWriter, responseData interface{}) {
	w.Header().Set("Content-Type", "applicaton/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(responseData)
}

func (r *JSONResponder) ErrorUnauthorized(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusUnauthorized, err)
}

func (r *JSONResponder) ErrorBadRequest(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusBadRequest, err)
}

func (r *JSONResponder) ErrorNotFound(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusNotFound, err)
}

func (r *JSONResponder) ErrorInternal(w http.ResponseWriter, err error) {
	r.sendError(w, http.StatusInternalServerError, err)
}

func (r *JSONResponder) sendError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}
