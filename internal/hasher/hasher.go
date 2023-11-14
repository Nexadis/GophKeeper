package hasher

import "golang.org/x/crypto/bcrypt"

type hasher struct {
	cost int
}

func New(cost int) hasher {
	return hasher{
		cost,
	}
}

func (h hasher) Hash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, h.cost)
}

func (h hasher) Compare(hash, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}
