package api

import (
	"encoding/json"
	"github.com/winsock/gochat/database"
	"net/http"
	"time"
)

type WebAPI struct {
	db *database.Database
}

type Error struct {
	Message string
	Time    time.Time
}

func Create(db *database.Database) *WebAPI {
	return &WebAPI{db: db}
}

func (api *WebAPI) writeJsonResponse(w http.ResponseWriter, value interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(value)
}

func (api *WebAPI) writeJsonError(w http.ResponseWriter, error string, statusCode int) error {
	return api.writeJsonResponse(w, Error{Message: error, Time: time.Now()}, statusCode)
}
