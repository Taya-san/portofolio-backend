package controller

import (
	"net/http"

	"github.com/Taya-san/portofolio-backend/config"
	"github.com/Taya-san/portofolio-backend/internal/controller/middleware"
	v1 "github.com/Taya-san/portofolio-backend/internal/controller/v1"
	projectusecase "github.com/Taya-san/portofolio-backend/internal/usecase/project"
	userusecase "github.com/Taya-san/portofolio-backend/internal/usecase/user"
	"github.com/Taya-san/portofolio-backend/pkg/jwt"
	"github.com/Taya-san/portofolio-backend/pkg/logger"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRouter(
	app *fiber.App,
	cfg *config.Config,
	u userusecase.User,
	p projectusecase.Project,
	jwtManager *jwt.Manager,
	l logger.Interface,
) {
	app.Use(middleware.Logger(l))
	app.Use(middleware.Recovery(l))

	if cfg.Metrics.Enabled {
		prometheus := fiberprometheus.New("portofolio-backend")
		prometheus.RegisterAt(app, "/metrics")
		app.Use(prometheus.Middleware)
	}

	if cfg.Swagger.Enabled {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	app.Get("/health", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

	apiV1Group := app.Group("/v1")
	{
		if cfg.Tracing.Enabled {
			apiV1Group.Use(otelfiber.Middleware())
		}

		v1.NewRoutes(
			apiV1Group,
			u,
			p,
			jwtManager,
			l,
		)
	}
}
