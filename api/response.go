package api

import (
	"github.com/winsock/gochat/database"
	"time"
)

type ErrorResponse struct {
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type SendMessageResponse struct {
	Message             database.Message `json:"message"`
	SentLiveToSender    bool             `json:"sentLiveToSender"`
	SentLiveToRecipient bool             `json:"sentLiveToRecipient"`
}
