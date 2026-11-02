package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Taya-san/portofolio-backend/config"
	"github.com/Taya-san/portofolio-backend/internal/controller"
	persistGithubProjectRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/github_project"
	persistProjectRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/project"
	persistUserRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/user"
	githubProjectApi "github.com/Taya-san/portofolio-backend/internal/repo/webapi/github_project"
	"github.com/Taya-san/portofolio-backend/internal/scheduler"
	githubproject "github.com/Taya-san/portofolio-backend/internal/usecase/github_project"
	"github.com/Taya-san/portofolio-backend/internal/usecase/project"
	"github.com/Taya-san/portofolio-backend/internal/usecase/user"
	"github.com/Taya-san/portofolio-backend/pkg/httpserver"
	"github.com/Taya-san/portofolio-backend/pkg/jwt"
	"github.com/Taya-san/portofolio-backend/pkg/logger"
	"github.com/Taya-san/portofolio-backend/pkg/tracing"
	"github.com/go-co-op/gocron/v2"
	_ "modernc.org/sqlite"
)

type useCases struct {
	project       project.Project
	user          user.User
	githubProject githubproject.GithubProject
}

func initUseCases(db *sql.DB, cfg *config.Config, jwtManager *jwt.Manager, h *http.Client) useCases {
	projectRepo := persistProjectRepo.New(db)
	userRepo := persistUserRepo.New(cfg)
	githubProjectRepo := persistGithubProjectRepo.New(db)
	githubProjectApi := githubProjectApi.New(h, cfg.Github.Token)

	return useCases{
		project:       project.New(projectRepo),
		user:          user.New(userRepo, jwtManager),
		githubProject: githubproject.New(githubProjectRepo, githubProjectApi),
	}
}

func initServer(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) *httpserver.Server {
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	controller.NewRouter(httpServer.App, cfg, uc.user, uc.project, uc.githubProject, jwtManager, l)

	return httpServer
}

func waitForShutdown(s *httpserver.Server, l logger.Interface, scheduler gocron.Scheduler) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error

	select {
	case sig := <-interrupt:
		l.Info("app - Run - signal: %s", sig.String())
	case err = <-s.Notify():
		l.Error(fmt.Errorf("app - Run - httpserver.Server.Notify: %w", err))
	}

	shutdownServer(s, l, scheduler)
}

func shutdownServer(s *httpserver.Server, l logger.Interface, scheduler gocron.Scheduler) {
	if err := s.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - Run - httpserver.Server.Shutdown: %w", err))
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := scheduler.ShutdownWithContext(ctx); err != nil {
		l.Error(fmt.Errorf("app - Run - gocron.Scheduler.ShutdownWithContext: %w", err))
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
	h := http.Client{
		Timeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	uc := initUseCases(db, cfg, jwtManager, &h)
	s := initServer(cfg, uc, jwtManager, l)
	scheduler, err := scheduler.NewScheduler(ctx, scheduler.GocronLogger{L: l}, uc.githubProject)
	if err != nil {
		l.Error(fmt.Errorf("app - Run - scheduler.NewScheduler: %w", err))
	}
	s.Start()
	waitForShutdown(s, l, scheduler)
}
