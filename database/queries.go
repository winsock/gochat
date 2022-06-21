package database

func (db *Database) InsertUser(user User) (User, error) {
	_, err := db.database.Exec("INSERT INTO users(uuid, username) VALUES (?, ?)", user.UUID, user.Username)
	return user, err
}

func (db *Database) FindUser(username string) (User, error) {
	var user User
	err := db.database.QueryRow("SELECT uuid, username FROM users WHERE username = ? LIMIT 1", username).Scan(&user.UUID, &user.Username)
	return user, err
}

func (db *Database) InsertMessage(message Message) (Message, error) {
	_, err := db.database.Exec("INSERT INTO messages(uuid, content, sender, recipient, created_at) VALUES (?, ?, ?, ?, ?)",
		message.Content, message.Sender.UUID, message.Recipient.UUID, message.CreatedAt)
	return message, err
}
