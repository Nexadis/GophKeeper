package users

import (
	"time"

	"github.com/Nexadis/GophKeeper/internal/models"
)

type user struct {
	id        int
	username  string
	hash      []byte
	createdAt time.Time
}

type User interface {
	ID() int
	Username() string
	Hash() []byte
	CreatedAt() time.Time
}

func (u user) ID() int {
	return u.id
}

func (u *user) SetID(id int) {
	u.id = id
}

func (u user) Username() string {
	return u.username
}

func (u user) Hash() []byte {
	return u.hash
}

func (u user) CreatedAt() time.Time {
	return u.createdAt
}

type UsersFactory struct {
	t models.TimeProvider
}

// New Создаёт фабрику пользователей. TimeProvider помогает создавать моки для тестирования создания пользователей с одним и тем же временем создания.
func New(t models.TimeProvider) *UsersFactory {
	return &UsersFactory{
		t,
	}
}

// New Создаёт нового пользователя, с заданными параметрами.
func (uc UsersFactory) New(username string, hash []byte) User {
	now := time.Now()
	if uc.t != nil {
		now = uc.t.Now()
	}
	return &user{
		username:  username,
		hash:      hash,
		createdAt: now,
	}
}
