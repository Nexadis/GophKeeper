package database

import (
	"context"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/models/datas"
	"github.com/Nexadis/GophKeeper/internal/models/users"
)

type PgDB struct {
	pg *sqlx.DB
	c  *config.DBConfig
}

func Connect(ctx context.Context, c *config.DBConfig) (*PgDB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", c.URI)
	if err != nil {
		return nil, err
	}
	return &PgDB{
		db,
		c,
	}, nil

}

func (pg *PgDB) Close() error {
	return pg.pg.Close()

}

func (pg *PgDB) Add(ctx context.Context, data datas.IData) error {
	return nil
}
func (pg *PgDB) GetByID(ctx context.Context, id int) (datas.IData, error) {
	return nil, nil
}
func (pg *PgDB) GetByUser(ctx context.Context, u users.User) ([]datas.IData, error) {
	return nil, nil
}
func (pg *PgDB) Update(ctx context.Context, data datas.IData) error {
	return nil
}
func (pg *PgDB) DeleteByID(ctx context.Context, id int) error {
	return nil
}
func (pg *PgDB) Ping(ctx context.Context) error {
	return pg.pg.PingContext(ctx)
}
func (pg *PgDB) GetUserByID(ctx context.Context, id int) (users.User, error) {
	return nil, nil
}
func (pg *PgDB) GetUserByName(ctx context.Context, username string) (users.User, error) {
	return nil, nil
}

func (pg *PgDB) AddUser(ctx context.Context, u users.User) error {
	return nil

}

func (pg *PgDB) DeleteUser(ctx context.Context, username string) error {
	return nil

}
