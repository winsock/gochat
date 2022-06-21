package database

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UUID     uuid.UUID
	Username string
}

type Message struct {
	UUID      uuid.UUID
	CreatedAt time.Time
	Content   string
	Sender    User
	Recipient User
}
