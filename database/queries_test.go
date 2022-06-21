package database

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDatabase_InsertUser(t *testing.T) {
	db := createTestDatabase(t)

	newUserUUID := uuid.New()
	user, err := db.InsertUser(User{
		UUID:     newUserUUID,
		Username: "testUser",
	})

	assert.Nil(t, err)
	assert.Equal(t, newUserUUID, user.UUID)
	assert.Equal(t, "testUser", user.Username)

}

func TestDatabase_FindUser(t *testing.T) {
	db := createTestDatabase(t)
	user := createTestUser(t, db, "test")

	foundUser, err := db.FindUser("test")

	assert.Nil(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user.UUID, foundUser.UUID)
	assert.Equal(t, user.Username, foundUser.Username)
}

func TestDatabase_InsertMessage(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	recipient := createTestUser(t, db, "test2")

	newMessageUuid := uuid.New()
	newMessageTime := time.Now()
	message, err := db.InsertMessage(Message{
		UUID:      newMessageUuid,
		CreatedAt: newMessageTime,
		Content:   "This is a test",
		Sender:    sender,
		Recipient: recipient,
	})

	assert.Nil(t, err)
	assert.Equal(t, newMessageUuid, message.UUID)
	assert.Equal(t, newMessageTime, message.CreatedAt)
	assert.Equal(t, "This is a test", message.Content)
	assert.Equal(t, sender.UUID, message.Sender.UUID)
	assert.Equal(t, sender.Username, message.Sender.Username)
	assert.Equal(t, recipient.UUID, message.Recipient.UUID)
	assert.Equal(t, recipient.Username, message.Recipient.Username)
}

func createTestDatabase(t *testing.T) *Database {
	db, err := Open()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	return db
}

func createTestUser(t *testing.T, db *Database, username string) User {
	user, err := db.InsertUser(User{
		UUID:     uuid.New(),
		Username: username,
	})
	assert.Nil(t, err)

	return user
}
