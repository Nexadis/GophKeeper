package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/Nexadis/GophKeeper/internal/app/httpserver"
	"github.com/Nexadis/GophKeeper/internal/app/services"
	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/database"
	"github.com/Nexadis/GophKeeper/internal/models/users"
)

type Server struct {
	config *config.AppConfig
	http   *httpserver.Server
}

func NewServer(c *config.AppConfig) Server {
	repo := database.New(c.DB.URI)
	uf := users.New(nil)
	as := services.NewAuth(repo, services.NewHash(), uf)
	ds := services.NewData(repo)
	s := Server{
		config: c,
	}
	if c.HTTP.Up {
		s.http = httpserver.New(c.HTTP, ds, as)
	}
	return s
}

func (s Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error { return s.http.Run(ctx) })
	return eg.Wait()
}
