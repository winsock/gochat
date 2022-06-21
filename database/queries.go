package database

func (db *Database) InsertUser(user User) (User, error) {
	err := db.database.QueryRow("INSERT INTO users(username) VALUES (?) RETURNING uuid", user.Username).Scan(&user.UUID)
	return user, err
}

func (db *Database) FindUser(username string) (User, error) {
	var user User
	err := db.database.QueryRow("SELECT uuid, username FROM users WHERE username = ? LIMIT 1", username).Scan(&user.UUID, &user.Username)
	return user, err
}

func (db *Database) InsertMessage(message Message) (Message, error) {
	err := db.database.QueryRow("INSERT INTO messages(content, sender, recipient, created_at) VALUES (?, ?, ?, ?) RETURNING uuid",
		message.Content, message.Sender.UUID, message.Recipient.UUID, message.CreatedAt).Scan(&message.UUID)
	return message, err
}
