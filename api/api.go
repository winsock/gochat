package api

import (
	"encoding/json"
	"github.com/winsock/gochat/database"
	"github.com/winsock/gochat/live"
	"log"
	"net/http"
	"strconv"
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
	if err != nil {
		log.Printf("API Error. %s, %s\n", message, err.Error())
	}
	return api.writeJsonResponse(w, ErrorResponse{Message: message, Time: time.Now()}, statusCode)
}

func (api *WebAPI) parseUint(formValue string, defaultValue uint64) (uint64, error) {
	if len(formValue) == 0 {
		return defaultValue, nil
	}
	return strconv.ParseUint(formValue, 10, 64)
}

func (api *WebAPI) parseTime(formValue string, defaultValue time.Time) (time.Time, error) {
	if len(formValue) == 0 {
		return defaultValue, nil
	}
	return time.Parse(time.RFC3339, formValue)
}
