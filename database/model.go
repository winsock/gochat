package database

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID     uuid.UUID `json:"uuid"`
	Username string    `json:"username"`
}

type Message struct {
	UUID      uuid.UUID `json:"uuid"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
	Sender    User      `json:"sender"`
	Recipient User      `json:"recipient"`
}
