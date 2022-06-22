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
	newMessageTime := time.Now().In(time.UTC)
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

func TestDatabase_FindMessagesForUser(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	recipient := createTestUser(t, db, "test2")
	message := createTestMessage(t, db, "test message", sender, recipient)

	foundMessages, err := db.FindMessagesForUser(recipient, message.CreatedAt, 0, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message.UUID, foundMessages[0].UUID)
	assert.Equal(t, message.Content, foundMessages[0].Content)
	assert.Equal(t, message.CreatedAt.Truncate(time.Nanosecond), foundMessages[0].CreatedAt.Truncate(time.Nanosecond))
	assert.Equal(t, message.Sender.UUID, foundMessages[0].Sender.UUID)
	assert.Equal(t, message.Sender.Username, foundMessages[0].Sender.Username)
	assert.Equal(t, message.Recipient.UUID, foundMessages[0].Recipient.UUID)
	assert.Equal(t, message.Recipient.Username, foundMessages[0].Recipient.Username)
}

func TestDatabase_FindMessagesForUserPagination(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	recipient := createTestUser(t, db, "test2")
	message := createTestMessage(t, db, "test message", sender, recipient)
	message2 := createTestMessage(t, db, "test message2", sender, recipient)

	foundMessages, err := db.FindMessagesForUser(recipient, message.CreatedAt, 0, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message.Content, foundMessages[0].Content)

	foundMessages, err = db.FindMessagesForUser(recipient, message.CreatedAt, 1, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message2.Content, foundMessages[0].Content)
}

func TestDatabase_FindMessagesForUserTooOld(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	recipient := createTestUser(t, db, "test2")
	_ = createTestMessage(t, db, "test message", sender, recipient)
	message2 := createTestMessage(t, db, "test message2", sender, recipient)

	foundMessages, err := db.FindMessagesForUser(recipient, message2.CreatedAt, 0, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message2.Content, foundMessages[0].Content)

	foundMessages, err = db.FindMessagesForUser(recipient, message2.CreatedAt, 1, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 0)
}

func TestDatabase_FindMessagesForUserMultipleSenders(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	sender2 := createTestUser(t, db, "test2")
	recipient := createTestUser(t, db, "test3")
	message := createTestMessage(t, db, "test message", sender, recipient)
	message2 := createTestMessage(t, db, "test message2", sender, recipient)
	message3 := createTestMessage(t, db, "test message3", sender2, recipient)

	foundMessages, err := db.FindMessagesForUser(recipient, message.CreatedAt, 0, 100)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 3)
	assert.Equal(t, message.Content, foundMessages[0].Content)
	assert.Equal(t, message2.Content, foundMessages[1].Content)
	assert.Equal(t, message3.Content, foundMessages[2].Content)
}

func TestDatabase_FindMessagesForUserFromSender(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	recipient := createTestUser(t, db, "test2")
	message := createTestMessage(t, db, "test message", sender, recipient)

	foundMessages, err := db.FindMessagesForUserFromSender(recipient, sender, message.CreatedAt, 0, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message.UUID, foundMessages[0].UUID)
	assert.Equal(t, message.Content, foundMessages[0].Content)
	assert.Equal(t, message.CreatedAt.Truncate(time.Nanosecond), foundMessages[0].CreatedAt.Truncate(time.Nanosecond))
	assert.Equal(t, message.Sender.UUID, foundMessages[0].Sender.UUID)
	assert.Equal(t, message.Sender.Username, foundMessages[0].Sender.Username)
	assert.Equal(t, message.Recipient.UUID, foundMessages[0].Recipient.UUID)
	assert.Equal(t, message.Recipient.Username, foundMessages[0].Recipient.Username)
}

func TestDatabase_FindMessagesForUserFromSenderPagination(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	recipient := createTestUser(t, db, "test2")
	message := createTestMessage(t, db, "test message", sender, recipient)
	message2 := createTestMessage(t, db, "test message2", sender, recipient)

	foundMessages, err := db.FindMessagesForUserFromSender(recipient, sender, message.CreatedAt, 0, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message.Content, foundMessages[0].Content)

	foundMessages, err = db.FindMessagesForUserFromSender(recipient, sender, message.CreatedAt, 1, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message2.Content, foundMessages[0].Content)
}

func TestDatabase_FindMessagesForUserFromSenderTooOld(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	recipient := createTestUser(t, db, "test2")
	_ = createTestMessage(t, db, "test message", sender, recipient)
	message2 := createTestMessage(t, db, "test message2", sender, recipient)

	foundMessages, err := db.FindMessagesForUserFromSender(recipient, sender, message2.CreatedAt, 0, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 1)
	assert.Equal(t, message2.Content, foundMessages[0].Content)

	foundMessages, err = db.FindMessagesForUserFromSender(recipient, sender, message2.CreatedAt, 1, 1)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 0)
}

func TestDatabase_FindMessagesForUserFromSenderMultipleSenders(t *testing.T) {
	db := createTestDatabase(t)
	sender := createTestUser(t, db, "test")
	sender2 := createTestUser(t, db, "test2")
	recipient := createTestUser(t, db, "test3")
	message := createTestMessage(t, db, "test message", sender, recipient)
	message2 := createTestMessage(t, db, "test message2", sender, recipient)
	_ = createTestMessage(t, db, "test message3", sender2, recipient)

	foundMessages, err := db.FindMessagesForUserFromSender(recipient, sender, message.CreatedAt, 0, 100)

	assert.Nil(t, err)
	assert.Len(t, foundMessages, 2)
	assert.Equal(t, message.Content, foundMessages[0].Content)
	assert.Equal(t, message2.Content, foundMessages[1].Content)
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

func createTestMessage(t *testing.T, db *Database, messageContents string, sender User, recipient User) Message {
	message, err := db.InsertMessage(Message{
		UUID:      uuid.New(),
		CreatedAt: time.Now().In(time.UTC),
		Content:   messageContents,
		Sender:    sender,
		Recipient: recipient,
	})
	assert.Nil(t, err)
	return message
}
