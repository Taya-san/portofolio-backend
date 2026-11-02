package v1

import (
	"net/http"

	"github.com/Taya-san/portofolio-backend/internal/controller/v1/request"
	"github.com/gofiber/fiber/v2"
)

// @Summary     Get githubProject
// @Description Get data from github_projects
// @ID          get-githubProject
// @Tags        githubProjects
// @Accept      json
// @Produce     json
// @Param       request body     request.GetByProjectUrl true "Project Url"
// @Success     201     {object} entity.GithubProject
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /edit/githubproject [get]
func (r *V1) getGithubProject(ctx *fiber.Ctx) error {
	var body request.GetByProjectUrl

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "REST API - V1 - getGithubProject")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "REST API - V1 - getGithubProject")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	project, err := r.g.GetByProjectUrl(ctx.UserContext(), body.ProjectUrl)
	if err != nil {
		r.l.Error(err, "REST API - V1 - getGithubProject")

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusCreated).JSON(project)
}
