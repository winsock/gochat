package database

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
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
