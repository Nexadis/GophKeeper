package database

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
	"github.com/Nexadis/GophKeeper/internal/models/users"
)

var ErrUserExist = errors.New("user is exist")
var ErrUserNotFound = errors.New("user not found")

type PgDB struct {
	db *pgxpool.Pool
	c  *config.DBConfig
}

func Connect(ctx context.Context, c *config.DBConfig) (*PgDB, error) {
	db, err := pgxpool.New(ctx, c.URI)
	if err != nil {
		return nil, err
	}
	return &PgDB{
		db,
		c,
	}, nil

}

func (pg *PgDB) Close() {
	pg.db.Close()
}

func (pg *PgDB) Add(ctx context.Context, data datas.IData) error {
	return nil
}
func (pg *PgDB) GetByID(ctx context.Context, id int) (datas.IData, error) {
	return nil, nil
}
func (pg *PgDB) GetByUser(ctx context.Context, u *users.User) ([]datas.IData, error) {
	return nil, nil
}
func (pg *PgDB) Update(ctx context.Context, data datas.IData) error {
	return nil
}
func (pg *PgDB) DeleteByID(ctx context.Context, id int) error {
	return nil
}
func (pg *PgDB) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}
func (pg *PgDB) GetUserByID(ctx context.Context, id int) (*users.User, error) {
	return nil, nil
}
func (pg *PgDB) GetUserByName(ctx context.Context, username string) (*users.User, error) {
	query := `SELECT id,username,hash,created_at FROM users WHERE username=@username`
	args := pgx.NamedArgs{
		"username": username,
	}
	logger.Debug(fmt.Sprintf("Get user %s from db", username))
	res := pg.db.QueryRow(ctx, query, args)
	var u users.User
	var hash string
	err := res.Scan(&u.ID, &u.Username, &hash, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NoDataFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("problem with db: %w", err)
	}
	h, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	u.Hash = append(u.Hash, h...)
	logger.Debug(fmt.Sprintf("Get user %s with id %d", u.Username, u.ID))
	return &u, nil
}

func (pg *PgDB) AddUser(ctx context.Context, u *users.User) error {
	query := `INSERT INTO users (username, hash, created_at) VALUES (@username, @hash, @created_at) RETURNING id`
	args := pgx.NamedArgs{
		"username":   u.Username,
		"hash":       hex.EncodeToString(u.Hash),
		"created_at": u.CreatedAt,
	}
	logger.Debug(fmt.Sprintf("Add user %s in db", u.Username))
	res := pg.db.QueryRow(ctx, query, args)
	var id int
	err := res.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {

			return ErrUserExist
		}
		return fmt.Errorf("problem with db: %w", err)
	}
	u.ID = id
	logger.Debug(fmt.Sprintf("Added user %s with id %d", u.Username, id))
	return nil

}

func (pg *PgDB) DeleteUser(ctx context.Context, username string) error {
	return nil

}
