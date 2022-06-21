package database

import (
	"database/sql"
	"time"
)

func (db *Database) InsertUser(user User) (User, error) {
	_, err := db.database.Exec("INSERT INTO users(uuid, username) VALUES (?, ?)", user.UUID, user.Username)
	return user, err
}

func (db *Database) FindUser(username string) (*User, error) {
	var user User
	err := db.database.QueryRow("SELECT uuid, username FROM users WHERE username = ? LIMIT 1", username).Scan(&user.UUID, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *Database) InsertMessage(message Message) (Message, error) {
	_, err := db.database.Exec("INSERT INTO messages(uuid, content, sender, recipient, created_at) VALUES (?, ?, ?, ?, ?)",
		message.UUID, message.Content, message.Sender.UUID, message.Recipient.UUID, message.CreatedAt.Format(time.RFC3339))
	return message, err
}

func (db *Database) FindMessagesForUser(user User, from time.Time, offset uint64, limit uint64) ([]Message, error) {
	rows, err := db.database.Query(`
	SELECT messages.uuid, messages.content, messages.created_at, senderUser.uuid, senderUser.username, recipientUser.uuid, recipientUser.username FROM messages
	LEFT JOIN users senderUser ON senderUser.uuid = messages.sender LEFT JOIN users recipientUser ON recipientUser.uuid = messages.recipient
	WHERE recipient = ? AND created_at >= ? ORDER BY messages.created_at LIMIT ? OFFSET ?
	`, user.UUID, from.Format(time.RFC3339), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseMessages(rows)
}

func (db *Database) FindMessagesForUserFromSender(user User, sender User, from time.Time, offset uint64, limit uint64) ([]Message, error) {
	rows, err := db.database.Query(`
	SELECT messages.uuid, messages.content, messages.created_at, senderUser.uuid, senderUser.username, recipientUser.uuid, recipientUser.username FROM messages
	LEFT JOIN users senderUser ON senderUser.uuid = messages.sender LEFT JOIN users recipientUser ON recipientUser.uuid = messages.recipient
	WHERE recipient = ? and sender = ? AND created_at >= ? ORDER BY messages.created_at LIMIT ? OFFSET ?
	`, user.UUID, sender.UUID, from.Format(time.RFC3339), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return parseMessages(rows)
}

func parseMessages(rows *sql.Rows) ([]Message, error) {
	messages := make([]Message, 0)
	for rows.Next() {
		var message Message
		var createdAtString string
		err := rows.Scan(&message.UUID, &message.Content, &createdAtString, &message.Sender.UUID, &message.Sender.Username, &message.Recipient.UUID, &message.Recipient.Username)
		if err != nil {
			return messages, err
		}
		message.CreatedAt, err = time.Parse(time.RFC3339, createdAtString)
		if err != nil {
			return messages, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}
