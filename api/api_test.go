package api

import (
	"encoding/json"
	"fmt"
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
	var response database.User
	decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Nil(t, decodeErr)
	assert.Nil(t, err)
	assert.Equal(t, *user, response)
}

func TestWebAPI_CreateDuplicateUser(t *testing.T) {
	api := createTestAPI(t)
	_ = createTestUser(t, api, "test")
	req := httptest.NewRequest(http.MethodGet, "/user/create?username=test", nil)
	w := httptest.NewRecorder()

	api.CreateUser(w, req)
	var response ErrorResponse
	decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	assert.Nil(t, decodeErr)
}

func TestWebAPI_SendMessage(t *testing.T) {
	api := createTestAPI(t)
	_ = createTestUser(t, api, "sender")
	recipient := createTestUser(t, api, "recipient")
	timeBeforeSend := time.Now()
	req := httptest.NewRequest(http.MethodGet, "/message/send?sender=sender&recipient=recipient&message=test", nil)
	w := httptest.NewRecorder()

	api.SendMessage(w, req)
	messages, err := api.db.FindMessagesForUser(recipient, timeBeforeSend, 0, 100)
	var response SendMessageResponse
	decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)

	assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	assert.Nil(t, decodeErr)
	assert.Nil(t, err)
	assert.Len(t, messages, 1)
	assert.Equal(t, messages[0], response.Message)
}

func TestWebAPI_SendMessage_InvalidUser(t *testing.T) {
	api := createTestAPI(t)
	req := httptest.NewRequest(http.MethodGet, "/message/send?sender=sender&recipient=recipient&message=test", nil)
	w := httptest.NewRecorder()

	api.SendMessage(w, req)
	var response ErrorResponse
	decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	assert.Nil(t, decodeErr)
}

func TestWebAPI_SearchMessages(t *testing.T) {
	api := createTestAPI(t)
	sender := createTestUser(t, api, "sender")
	recipient := createTestUser(t, api, "recipient")
	message := createTestMessage(t, api, "test message", sender, recipient)
	req := httptest.NewRequest(http.MethodGet, "/message/search?recipient=recipient", nil)
	w := httptest.NewRecorder()

	api.SearchMessages(w, req)
	messages, err := api.db.FindMessagesForUser(recipient, message.CreatedAt, 0, 100)
	var response SearchResponse
	decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Nil(t, decodeErr)
	assert.Nil(t, err)
	assert.Equal(t, messages, response.Messages)
	assert.Equal(t, uint64(1), response.Count)
	assert.Equal(t, uint64(0), response.Offset)
	assert.Equal(t, uint64(100), response.Limit)
}

func TestWebAPI_SearchMessagesTooOld(t *testing.T) {
	api := createTestAPI(t)
	sender := createTestUser(t, api, "sender")
	recipient := createTestUser(t, api, "recipient")
	message := createTestMessage(t, api, "test message", sender, recipient)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/message/search?recipient=recipient&from=%s", message.CreatedAt.Add(time.Minute).Format(time.RFC3339)), nil)
	w := httptest.NewRecorder()

	api.SearchMessages(w, req)
	var response SearchResponse
	decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)

	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	assert.Nil(t, decodeErr)
	assert.Len(t, response.Messages, 0)
	assert.Equal(t, response.Count, uint64(0))
}

func TestWebAPI_SearchMessages_InvalidCalls(t *testing.T) {
	api := createTestAPI(t)
	_ = createTestUser(t, api, "sender")
	_ = createTestUser(t, api, "recipient")

	t.Run("No Users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/message/search", nil)
		w := httptest.NewRecorder()
		api.SearchMessages(w, req)
		var response ErrorResponse
		decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, decodeErr)
	})
	t.Run("No Such Recipient", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/message/search?recipient=test", nil)
		w := httptest.NewRecorder()
		api.SearchMessages(w, req)
		var response ErrorResponse
		decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Nil(t, decodeErr)
	})
	t.Run("No Such Sender", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/message/search?recipient=recipient&sender=test", nil)
		w := httptest.NewRecorder()
		api.SearchMessages(w, req)
		var response ErrorResponse
		decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Nil(t, decodeErr)
	})
	t.Run("Invalid Date", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/message/search?recipient=recipient&from=baddate", nil)
		w := httptest.NewRecorder()
		api.SearchMessages(w, req)
		var response ErrorResponse
		decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, decodeErr)
	})
	t.Run("Invalid Offset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/message/search?recipient=recipient&offset=notanumber", nil)
		w := httptest.NewRecorder()
		api.SearchMessages(w, req)
		var response ErrorResponse
		decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, decodeErr)
	})
	t.Run("Invalid Limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/message/search?recipient=recipient&limit=notanumber", nil)
		w := httptest.NewRecorder()
		api.SearchMessages(w, req)
		var response ErrorResponse
		decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, decodeErr)
	})
	t.Run("Invalid Limit Too Large", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/message/search?recipient=recipient&limit=1001", nil)
		w := httptest.NewRecorder()
		api.SearchMessages(w, req)
		var response ErrorResponse
		decodeErr := json.NewDecoder(w.Result().Body).Decode(&response)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Nil(t, decodeErr)
	})
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
