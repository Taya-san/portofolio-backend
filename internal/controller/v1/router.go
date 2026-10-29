package v1

import (
	"github.com/Taya-san/portofolio-backend/internal/controller/middleware"
	projectusecase "github.com/Taya-san/portofolio-backend/internal/usecase/project"
	userusecase "github.com/Taya-san/portofolio-backend/internal/usecase/user"
	"github.com/Taya-san/portofolio-backend/pkg/jwt"
	"github.com/Taya-san/portofolio-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewRoutes(
	apiV1Group fiber.Router,
	u userusecase.User,
	p projectusecase.Project,
	jwtManager *jwt.Manager,
	l logger.Interface,
) {
	r := &V1{
		u: u,
		p: p,
		l: l,
		v: validator.New(validator.WithRequiredStructEnabled()),
	}

	// Public routes
	authGroup := apiV1Group.Group("/auth")
	{
		authGroup.Post("/login", r.login)
	}

	apiV1Group.Get("/projects", r.listProjects)
	apiV1Group.Get("/project/:id", r.getProject)

	protectedEdit := apiV1Group.Group("/edit", middleware.Auth(jwtManager))
	{
		protectedEdit.Post("/", r.createProject)
		protectedEdit.Put("/:id", r.updateProject)
		protectedEdit.Delete("/:id", r.deleteProject)
	}
}
