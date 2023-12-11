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
var ErrDataNotFound = errors.New("data not found")

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

func (pg *PgDB) Add(ctx context.Context, dlist []datas.Data) error {
	query := `INSERT INTO datas (
	user_id,
	dtype,
	description,
	value,
	created_at,
	edited_at) 

	VALUES (
  @user_id,
  @dtype,
  @description,
  @value ,
  @created_at,
  @edited_at) `
	batch := &pgx.Batch{}
	for _, data := range dlist {
		args := pgx.NamedArgs{
			"user_id":     data.UserID,
			"dtype":       data.Type,
			"description": data.Description,
			"value":       data.Value,
			"created_at":  data.CreatedAt,
			"edited_at":   data.EditedAt,
		}
		batch.Queue(query, args)
	}
	results := pg.db.SendBatch(ctx, batch)
	defer results.Close()
	for _, d := range dlist {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("can't add data '%s': %w", d.Value, err)
		}

	}
	uid := dlist[0].UserID

	logger.Debug(fmt.Sprintf("Added batch with data from uid=%d", uid))
	return results.Close()

}
func (pg *PgDB) GetByID(ctx context.Context, id int) (*datas.Data, error) {
	query := `SELECT id,user_id,dtype,description, value, created_at, edited_at FROM datas WHERE id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}
	logger.Debug("Get data with id %d from db", id)
	res := pg.db.QueryRow(ctx, query, args)
	var d datas.Data
	err := res.Scan(&d.ID, &d.UserID, &d.Type, &d.Description, &d.Value, &d.CreatedAt, &d.EditedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NoDataFound {
			return nil, ErrDataNotFound
		}
	}
	logger.Debug("Get data: %v", d)
	return &d, nil

}
func (pg *PgDB) GetByUser(ctx context.Context, uid int) ([]datas.Data, error) {
	query := `SELECT * FROM datas WHERE user_id=@user_id`
	args := pgx.NamedArgs{
		"user_id": uid,
	}
	rows, err := pg.db.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("unable to query data of uid=%d : %w", uid, err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows,
		func(row pgx.CollectableRow) (datas.Data, error) {
			var d datas.Data
			err := row.Scan(
				&d.ID,
				&d.UserID,
				&d.Type,
				&d.Description,
				&d.Value,
				&d.CreatedAt,
				&d.EditedAt,
			)
			return d, err
		},
	)

}
func (pg *PgDB) Update(ctx context.Context, dlist []datas.Data) error {
	query := "UPDATE datas SET value=@value, edited_at=@edited_at, description=@description WHERE id=@id AND user_id=@user_id"
	batch := &pgx.Batch{}
	for _, data := range dlist {
		args := pgx.NamedArgs{
			"id":          data.ID,
			"user_id":     data.UserID,
			"value":       data.Value,
			"edited_at":   data.EditedAt,
			"description": data.Description,
		}
		batch.Queue(query, args)
	}
	results := pg.db.SendBatch(ctx, batch)
	defer results.Close()
	for _, d := range dlist {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.NoDataFound {
				return fmt.Errorf("can't update data id=%d: %w", d.ID, ErrDataNotFound)

			}
			return fmt.Errorf("can't update data id=%d: %w", d.ID, err)

		}

	}
	return results.Close()
}
func (pg *PgDB) DeleteByIDs(ctx context.Context, uid int, ids []int) error {
	query := `DELETE FROM datas WHERE id=@id AND user_id=@uid`
	batch := &pgx.Batch{}
	for _, id := range ids {
		args := pgx.NamedArgs{
			"id":  id,
			"uid": uid,
		}
		batch.Queue(query, args)
	}
	results := pg.db.SendBatch(ctx, batch)
	defer results.Close()
	for _, id := range ids {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("can't delete data id=%d : %w", id, err)
		}

	}

	logger.Debug(fmt.Sprintf("Data deleted by request of uid=%d", uid))
	return results.Close()
}
func (pg *PgDB) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
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
		return nil, fmt.Errorf("problem with db in get user by name for '%s' : %w", username, err)
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
		return fmt.Errorf("problem with db in add user '%s' : %w", u.Username, err)
	}
	u.ID = id
	logger.Debug(fmt.Sprintf("Added user %s with id %d", u.Username, id))
	return nil

}

func (pg *PgDB) DeleteUser(ctx context.Context, username string) error {
	return nil

}
