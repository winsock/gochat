package api

import (
	"github.com/google/uuid"
	"github.com/winsock/gochat/database"
	"log"
	"net/http"
)

func (api *WebAPI) CreateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	newUser, err := api.db.InsertUser(database.User{UUID: uuid.New(), Username: username})
	if err != nil {
		_ = api.writeJsonError(w, "Error creating user", http.StatusInternalServerError, err)
		return
	}

	log.Printf("New user created %+v\n", newUser)
	_ = api.writeJsonResponse(w, newUser, http.StatusCreated)
}
