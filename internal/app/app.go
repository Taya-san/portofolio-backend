package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Taya-san/portofolio-backend/config"
	"github.com/Taya-san/portofolio-backend/internal/controller"
	"github.com/Taya-san/portofolio-backend/internal/migrate"
	persistGithubProjectRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/github_project"
	persistProjectRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/project"
	persistUserRepo "github.com/Taya-san/portofolio-backend/internal/repo/persist/user"
	webGithubProjectApi "github.com/Taya-san/portofolio-backend/internal/repo/webapi/github_project"
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
	githubProjectApi := webGithubProjectApi.New(h, cfg.Github.Token)

	return useCases{
		project:       project.New(projectRepo, githubProjectRepo),
		user:          user.New(userRepo, jwtManager),
		githubProject: githubproject.New(githubProjectRepo, githubProjectApi),
	}
}

func initServer(cfg *config.Config, uc useCases, jwtManager *jwt.Manager, l logger.Interface) *httpserver.Server {
	httpServer := httpserver.New(l, httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
	controller.NewRouter(httpServer.App, cfg, uc.user, uc.project, uc.githubProject, jwtManager, l)

	return httpServer
}

func waitForShutdown(ctx context.Context, s *httpserver.Server, l logger.Interface, sch gocron.Scheduler) {
	select {
	case <-ctx.Done():
		l.Info("app - waitForShutdown - kill/term signal received")
	case err := <-s.Notify():
		l.Error(fmt.Errorf("app - waitForShutdown - httpserver.Server.Notify: %w", err))
	}

	shutdownServer(s, l, sch)
}

func shutdownServer(s *httpserver.Server, l logger.Interface, sch gocron.Scheduler) {
	if err := s.Shutdown(); err != nil {
		l.Error(fmt.Errorf("app - shutdownServer - httpserver.Server.Shutdown: %w", err))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := sch.ShutdownWithContext(ctx); err != nil {
		l.Error(fmt.Errorf("app - shutdownServer - gocron.Scheduler.ShutdownWithContext: %w", err))
	}
}

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	ctx := context.Background()

	db, err := sql.Open("sqlite",
		"file:data.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)",
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - sql.Open: %w", err))
	}

	db.SetMaxOpenConns(1)
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		l.Fatal(fmt.Errorf("app - Run - db.PingContext: %w", err))
	}

	err = migrate.DB(ctx, db)
	if err != nil {
		db.Close()
		l.Fatal(fmt.Errorf("app - Run - migrate.DB: %w", err))
	}

	args := os.Args[1:]

	if len(args) > 1 && !strings.HasPrefix(args[0], "-") {
		cmd := strings.Join(args, " ")
		switch cmd {
		case "serve":
			runServer(ctx, cfg, db, l)
		case "refresh mirror":
			refreshMirror(ctx, cfg, db, l)
		default:
			printGlobalHelp()
		}
	}
}

func printGlobalHelp() {
	fmt.Println("Usage: backend <command>")
	fmt.Println("\nAvailable Commands:")
	fmt.Println("		serve								start the portofolio-website backend server")
	fmt.Println("		refresh mirror			refresh the github projects mirror")
}

func refreshMirror(ctx context.Context, cfg *config.Config, db *sql.DB, l *logger.Logger) {
	githubProjectRepo := persistGithubProjectRepo.New(db)
	githubProjectApi := webGithubProjectApi.New(&http.Client{Timeout: 5 * time.Second}, cfg.Github.Token)

	githubProjectUc := githubproject.New(githubProjectRepo, githubProjectApi)

	err := githubProjectUc.RefreshMirror(ctx)
	if err != nil {
		db.Close()
		l.Fatal(fmt.Errorf("app - refreshMirror - githubProjectUc.RefreshMirror: %w", err))
	}
}

func runServer(ctx context.Context, cfg *config.Config, db *sql.DB, l *logger.Logger) {
	shutdownTracing, err := tracing.New(ctx, tracing.Config{
		Enabled:     cfg.Tracing.Enabled,
		ServiceName: cfg.App.Name,
		Version:     cfg.App.Version,
		Prettier:    cfg.Tracing.Prettier,
		SampleRate:  cfg.Tracing.SampleRate,
	})
	if err != nil {
		db.Close()
		l.Fatal(fmt.Errorf("app - runServer - tracing.New: %w", err))
	}
	defer func() {
		if err := shutdownTracing(ctx); err != nil {
			l.Error(fmt.Errorf("app - runServer - shutdownTracing: %w", err))
		}
	}()

	jwtManager := jwt.New(cfg.JWT.Secret, cfg.JWT.TokenExpiry)
	h := http.Client{
		Timeout: 5 * time.Second,
	}

	shutdownCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	uc := initUseCases(db, cfg, jwtManager, &h)
	s := initServer(cfg, uc, jwtManager, l)

	sch, err := scheduler.NewScheduler(shutdownCtx, scheduler.GocronLogger{L: l}, uc.githubProject)
	if err != nil {
		db.Close()
		l.Fatal(fmt.Errorf("app - runServer - scheduler.NewScheduler: %w", err))
	}

	s.Start()
	sch.Start()
	waitForShutdown(shutdownCtx, s, l, sch)
}
