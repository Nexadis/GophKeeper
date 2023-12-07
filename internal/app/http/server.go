package http

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/Nexadis/GophKeeper/internal/app/services"
	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
)

type Server struct {
	config      *config.HTTPServerConfig
	e           *echo.Echo
	dataService *services.Data
	authService *services.Auth
}

func New(c *config.HTTPServerConfig, d *services.Data, a *services.Auth) *Server {
	e := echo.New()
	hs := &Server{
		c,
		e,
		d,
		a,
	}
	hs.mountHandlers()
	return hs
}

func (hs *Server) mountHandlers() {
	hs.e.Use(middleware.Logger(),
		middleware.Gzip(),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:     true,
			LogStatus:  true,
			LogMethod:  true,
			LogLatency: true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				logger.Info(
					fmt.Sprintf(
						"request: %s %d URI=%s latency=%v",
						v.Method,
						v.Status,
						v.URI,
						v.Latency,
					),
				)
				return nil
			},
		},
		),
	)
	hs.e.Static("/", hs.config.FrontDir)
	hs.e.GET(APIDownload, hs.Download)
	hs.e.POST(APIRegister, hs.Register)
	hs.e.POST(APILogin, hs.Login)

	hs.mustSecret()

	apiv1 := hs.e.Group(APIv1)
	{
		config := echojwt.Config{
			NewClaimsFunc: func(c echo.Context) jwt.Claims {
				return new(Claims)
			},
			SigningKey: hs.config.JWTSecret,
		}
		apiv1.Use(echojwt.WithConfig(config))
		apiv1.GET(APIPing, hs.Ping)

		data := apiv1.Group(APIData)
		{
			data.GET("", hs.GetData)
			data.POST("", hs.PostData)
			data.DELETE("", hs.DeleteData)
			data.PATCH("", hs.UpdateData)
		}

	}

}

func (hs *Server) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		var err error
		defer cancel()
		if hs.config.TLS {
			err = hs.e.StartTLS(hs.config.Address, hs.config.CrtFile, hs.config.KeyFile)
		} else {
			logger.Info("Server is running without TLS. Be careful, your password may be intercepted!")
			err = hs.e.Start(hs.config.Address)
		}
		if err != nil {
			logger.Error(err)
		}

	}()

	<-ctx.Done()
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return hs.e.Shutdown(ctx)

}

func (hs *Server) mustSecret() {
	if hs.config.JWTSecret == nil {
		buf := make([]byte, 32)
		_, err := rand.Read(buf)
		if err != nil {
			log.Fatal(err)
		}
		hs.config.JWTSecret = buf
	}

}
