package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/Nexadis/GophKeeper/internal/app"
	"github.com/Nexadis/GophKeeper/internal/config"
	"github.com/Nexadis/GophKeeper/internal/logger"
)

func main() {
	c := config.MustClientConfig()
	logger.Init(c.Log)
	a := app.NewClient(c)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer cancel()

	err := a.Run(ctx)
	if err != nil {
		logger.Error(err)
	}

}
