package api

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/winsock/gochat/database"
	"github.com/winsock/gochat/live"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebAPI_CreateUser(t *testing.T) {
	api := createTestAPI(t)

	req := httptest.NewRequest(http.MethodGet, "/user/create?username=test", nil)
	w := httptest.NewRecorder()
	api.CreateUser(w, req)
	user, err := api.db.FindUser("test")

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Nil(t, err)
	assert.Equal(t, "test", user.Username)
}

func createTestAPI(t *testing.T) *WebAPI {
	db, err := database.Open()
	assert.Nil(t, err)
	assert.NotNil(t, db)
	// TODO, ideally I would create a common interface with all of the database methods that I would then implement per test
	// This would allow testing the API without touching any of the database code, not implemented in this demo due to time constraints
	return Create(db, live.Create(db))
}

func createTestUser(t *testing.T, api *WebAPI, username string) database.User {
	user, err := api.db.InsertUser(database.User{
		UUID:     uuid.New(),
		Username: username,
	})
	assert.Nil(t, err)

	return user
}

func createTestMessage(t *testing.T, api *WebAPI, messageContents string, sender database.User, recipient database.User) database.Message {
	message, err := api.db.InsertMessage(database.Message{
		UUID:      uuid.New(),
		CreatedAt: time.Now(),
		Content:   messageContents,
		Sender:    sender,
		Recipient: recipient,
	})
	assert.Nil(t, err)
	return message
}
