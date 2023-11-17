package users

import (
	"time"
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

func New(username string, hash []byte) *user {
	now := time.Now()
	return &user{
		username:  username,
		hash:      hash,
		createdAt: now,
	}
}

func (u user) ID() int {
	return u.id
}

func (u *user) SetID(id int) {
	u.id = id
	u.createdAt = time.Now()
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
