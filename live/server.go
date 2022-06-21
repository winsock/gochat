package live

import (
	"context"
	"github.com/google/uuid"
	"github.com/winsock/gochat/database"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	db          *database.Database
	ctx         context.Context
	clientMutex sync.RWMutex
	clients     map[uuid.UUID]*Client
}

func (server *Server) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	user, err := server.db.FindUser(r.FormValue("username"))
	if err != nil {
		log.Printf("Error locating user %s, %s\n", r.FormValue("username"), err.Error())
		return
	}

	// Check if we already have a client connected for the user, if so refuse connection
	server.clientMutex.RLock()
	if client, ok := server.clients[user.UUID]; ok {
		log.Printf("%s tried to connect with username %s but %s is already connected with that username", r.RemoteAddr, user.Username, client.remoteAddr)
		http.Error(w, "Client already connected with the same username!", http.StatusConflict)
		return
	}
	server.clientMutex.RUnlock()

	// Create the client
	client, err := Open(user, server.ctx, w, r)
	if err != nil {
		log.Printf("Error creating websocket client! %s\n", err.Error())
		http.Error(w, "An internal error has occurred", http.StatusInternalServerError)
	}

	server.clientMutex.Lock()
	server.clients[user.UUID] = client
	server.clientMutex.Unlock()
}
