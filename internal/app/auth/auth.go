package auth

import (
	"context"

	"github.com/Nexadis/GophKeeper/internal/users"
)

type Hasher interface {
	Compare(hash, password []byte) error
	Hash(password []byte) ([]byte, error)
}

type auth struct {
	hasher Hasher
}

func New(h Hasher) auth {
	return auth{
		h,
	}
}

func (a auth) Hash(password string) ([]byte, error) {
	h, err := a.hasher.Hash([]byte(password))
	return h, err
}

func (a auth) Auth(ctx context.Context, u users.User, password string) error {
	return a.hasher.Compare(u.Hash(), []byte(password))
}
