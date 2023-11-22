package users

import (
	"time"

	"github.com/Nexadis/GophKeeper/internal/models"
)

type User struct {
	ID        int
	Username  string
	Hash      []byte
	CreatedAt time.Time
}

type UsersFactory struct {
	t models.TimeProvider
}

// New Создаёт фабрику пользователей. TimeProvider помогает создавать моки для тестирования создания пользователей с одним и тем же временем создания.
func New(username string, hash []byte) User {
	now := time.Now()
	return User{
		Username:  username,
		Hash:      hash,
		CreatedAt: now,
	}
}
