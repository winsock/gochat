package main

import (
	"github.com/winsock/gochat/api"
	"github.com/winsock/gochat/database"
	"github.com/winsock/gochat/live"
	"log"
	"net/http"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Printf("Error opening database! %s\n", err.Error())
	}
	liveServer := live.Create(db)
	restAPI := api.Create(db, liveServer)

	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/live", liveServer.ServeWebsocket)
	httpMux.HandleFunc("/user/create", restAPI.CreateUser)
	httpMux.HandleFunc("/message/send", restAPI.SendMessage)
	//httpMux.HandleFunc("/message/search", restAPI.SearchMessages)

	err = http.ListenAndServe(":8080", httpMux)
	log.Fatalf("Error running HTTP server %s\n", err.Error())
}
