package v1

import (
	projectusecase "github.com/Taya-san/portofolio-backend/internal/usecase/project"
	userusecase "github.com/Taya-san/portofolio-backend/internal/usecase/user"
	"github.com/Taya-san/portofolio-backend/pkg/logger"
	"github.com/go-playground/validator/v10"
)

type V1 struct {
	u userusecase.User
	p projectusecase.Project
	l logger.Interface
	v *validator.Validate
}
