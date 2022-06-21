package api

import (
	"github.com/winsock/gochat/database"
	"net/http"
)

func (api *WebAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	newUser, err := api.db.InsertUser(database.User{Username: username})
	if err != nil {
		_ = api.writeJsonError(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	_ = api.writeJsonResponse(w, newUser, http.StatusOK)
}
