package users

import (
	"time"
)

type User struct {
	ID        int
	Username  string
	Hash      []byte
	CreatedAt time.Time
}

// New Создаёт пользователя с заданными полями
func New(username string, hash []byte) User {
	now := time.Now()
	return User{
		Username:  username,
		Hash:      hash,
		CreatedAt: now,
	}
}
