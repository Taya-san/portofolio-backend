package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Taya-san/portofolio-backend/config"
	"github.com/Taya-san/portofolio-backend/internal/controller"
	persistProjectRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/project"
	persistUserRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/user"
	"github.com/Taya-san/portofolio-backend/internal/usecase/project"
	"github.com/Taya-san/portofolio-backend/internal/usecase/user"
	"github.com/Taya-san/portofolio-backend/pkg/httpserver"
	"github.com/Taya-san/portofolio-backend/pkg/jwt"
	"github.com/Taya-san/portofolio-backend/pkg/logger"
	"github.com/Taya-san/portofolio-backend/pkg/tracing"
	_ "modernc.org/sqlite"
)

type useCases struct {
	project project.Project
	user    user.User
}

func initUseCases(db *sql.DB, cfg *config.Config, jwtManager *jwt.Manager) useCases {
	projectRepo := persistProjectRepo.New(db)
	userRepo := persistUserRepo.New(cfg)

	return useCases{
		project: project.New(projectRepo),
		user:    user.New(userRepo, jwtManager),
	}
}

func initServer(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) *httpserver.Server {
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	controller.NewRouter(httpServer.App, cfg, uc.user, uc.project, jwtManager, l)

	return httpServer
}

func waitForShutdown(s *httpserver.Server, l logger.Interface) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.Notify():
		l.Error(fmt.Errorf("app - Run - httpserver.Server.Notify: %w", err))
	}

	shutdownServer(s, l)
}

func shutdownServer(s *httpserver.Server, l logger.Interface) {
	if err := s.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpserver.Server.Shutdown: %w", err))
	}
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	ctx := context.Background()

	shutdownTracing, err := tracing.New(ctx, tracing.Config{
		Enabled:     cfg.Tracing.Enabled,
		ServiceName: cfg.App.Name,
		Version:     cfg.App.Version,
		Prettier:    cfg.Tracing.Prettier,
		SampleRate:  cfg.Tracing.SampleRate,
	})
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - tracing.New: %w", err))
	}
	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			l.Error(fmt.Errorf("app - Run - shutdownTracing: %w", err))
		}
	}()

	db, err := sql.Open("sqlite",
		"file:data.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)",
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - sql.Open: %w", err))
	}
	db.SetMaxOpenConns(1)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		l.Fatal(fmt.Errorf("app - Run - db.PingContext: %w", err))
	}

	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)

	uc := initUseCases(db, cfg, jwtManager)
	s := initServer(cfg, uc, jwtManager, l)
	s.Start()
	waitForShutdown(s, l)
}
