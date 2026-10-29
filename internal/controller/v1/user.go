package v1

import (
	"errors"
	"net/http"

	"github.com/Taya-san/portofolio-backend/internal/controller/v1/request"
	"github.com/Taya-san/portofolio-backend/internal/controller/v1/response"
	"github.com/Taya-san/portofolio-backend/internal/entity"
	userusecase "github.com/Taya-san/portofolio-backend/internal/usecase/user"
	"github.com/gofiber/fiber/v2"
)

// @Summary     Login
// @Description Authenticate user and get JWT token
// @ID          login
// @Tags        auth
// @Accept      json
// @Produce     json
// @Param       request body     request.Login true "Login credentials"
// @Success     200     {object} response.Token
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Router      /auth/login [post]
func (r *V1) login(ctx *fiber.Ctx) error {
	var body request.Login

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "REST API - V1 - login")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "REST API - V1 - login")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	token, err := r.u.Login(ctx.UserContext(), userusecase.UserItem{Password: body.Password})
	if err != nil {
		r.l.Error(err, "REST API - V1 - login")

		if errors.Is(err, entity.ErrInvalidCredentials) {
			return errorResponse(ctx, http.StatusUnauthorized, "invalid credentials")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(response.Token{Token: token})
}
