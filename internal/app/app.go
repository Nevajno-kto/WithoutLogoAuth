// Package app configures and runs application.
package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/nevajno-kto/without-logo-auth/config"
	"github.com/nevajno-kto/without-logo-auth/internal/controller/rest"
	"github.com/nevajno-kto/without-logo-auth/internal/usecase"
	"github.com/nevajno-kto/without-logo-auth/internal/usecase/repo/psql"
	"github.com/nevajno-kto/without-logo-auth/pkg/httpserver"
	"github.com/nevajno-kto/without-logo-auth/pkg/logger"
	"github.com/nevajno-kto/without-logo-auth/pkg/postgres"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	//pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(2))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Use case
	authUseCase := usecase.NewAuth(
		psql.NewAuthRepo(pg),
		psql.NewClientsRepo(pg),
		psql.NewPemissionsRepo(pg),
	)

	// HTTP Server
	handler := gin.New()
	rest.NewRouter(handler, l, authUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

}
