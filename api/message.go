package api

import (
	"errors"
	"github.com/google/uuid"
	"github.com/winsock/gochat/database"
	"github.com/winsock/gochat/live"
	"log"
	"net/http"
	"time"
)

func (api *WebAPI) SendMessage(w http.ResponseWriter, r *http.Request) {
	sender, err := api.db.FindUser(r.FormValue("sender"))
	if err != nil {
		_ = api.writeJsonError(w, "Unable to find sender with the provided username", http.StatusNotFound, err)
		return
	}
	recipient, err := api.db.FindUser(r.FormValue("recipient"))
	if err != nil {
		_ = api.writeJsonError(w, "Unable to find recipient with the provided username", http.StatusNotFound, err)
		return
	}
	message := r.FormValue("message")
	if len(message) == 0 {
		_ = api.writeJsonError(w, "Message cannot be empty", http.StatusBadRequest, err)
		return
	}

	newMessage, err := api.db.InsertMessage(database.Message{
		UUID:      uuid.New(),
		CreatedAt: time.Now(),
		Content:   message,
		Sender:    sender,
		Recipient: recipient,
	})
	if err != nil {
		_ = api.writeJsonError(w, "Error while sending message", http.StatusInternalServerError, err)
		return
	}

	sentLiveToSender := false
	sentLiveToRecipient := false

	// Send live to the sender if they are connected
	err = api.liveServer.RunWithClient(sender.UUID, func(client *live.Client) error {
		return client.SendMessage(newMessage)
	})
	// Ignore no such client errors, clients may not be always connected
	if err != nil && !errors.Is(err, live.NoSuchClientErr) {
		_ = api.writeJsonError(w, "Error while sending message to live client of the sender", http.StatusInternalServerError, err)
		return
	}
	sentLiveToSender = err == nil

	// Send live to the recipient if they are connected
	err = api.liveServer.RunWithClient(recipient.UUID, func(client *live.Client) error {
		return client.SendMessage(newMessage)
	})
	// Ignore no such client errors, clients may not be always connected
	if err != nil && !errors.Is(err, live.NoSuchClientErr) {
		_ = api.writeJsonError(w, "Error while sending message to live client of the recipient", http.StatusInternalServerError, err)
		return
	}
	sentLiveToRecipient = err == nil

	response := SendMessageResponse{
		Message:             newMessage,
		SentLiveToSender:    sentLiveToSender,
		SentLiveToRecipient: sentLiveToRecipient,
	}
	log.Printf("Message sent %+v\n", response)
	_ = api.writeJsonResponse(w, response, http.StatusOK)
}
