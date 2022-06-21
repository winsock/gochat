package database

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       uuid.UUID
	Username string
}

type Message struct {
	ID        uuid.UUID
	CreatedAt time.Time
	Contents  string
	Sender    uuid.UUID
	Recipient uuid.UUID
}
