package app

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/Nexadis/GophKeeper/internal/app/http"
	"github.com/Nexadis/GophKeeper/internal/app/services"
	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/database"
)

// Server - структура для работы с сервером
type Server struct {
	config *config.ServerConfig
	http   *http.Server
}

// NewServer - создаёт сервер с заданным конфигом
func NewServer(c *config.ServerConfig) (*Server, error) {
	s := Server{
		config: c,
	}
	return &s, nil
}

// Run - запускает сервер
func (s Server) Run(ctx context.Context) error {
	time.Sleep(s.config.WarmUp)
	dbctx, cancel := context.WithTimeout(
		ctx,
		time.Duration(s.config.DB.Timeout)*time.Second,
	)
	defer cancel()
	repo, err := database.Connect(dbctx, s.config.DB)
	if err != nil {
		return err
	}
	as := services.NewAuth(repo, services.NewHash())
	ds := services.NewData(repo)
	if s.config.HTTP.Up {
		s.http = http.New(s.config.HTTP, ds, as)
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return s.http.Run(ctx) })
	defer repo.Close()
	return eg.Wait()
}
