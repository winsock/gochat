package api

import (
	"encoding/json"
	"github.com/winsock/gochat/database"
	"github.com/winsock/gochat/live"
	"log"
	"net/http"
	"time"
)

type WebAPI struct {
	db         *database.Database
	liveServer *live.Server
}

func Create(db *database.Database, liveServer *live.Server) *WebAPI {
	return &WebAPI{db: db, liveServer: liveServer}
}

func (api *WebAPI) writeJsonResponse(w http.ResponseWriter, value interface{}, statusCode int) error {
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(value)
}

func (api *WebAPI) writeJsonError(w http.ResponseWriter, message string, statusCode int, err error) error {
	log.Printf("API Error. %s, %s\n", message, err.Error())
	return api.writeJsonResponse(w, ErrorResponse{Message: message, Time: time.Now()}, statusCode)
}
