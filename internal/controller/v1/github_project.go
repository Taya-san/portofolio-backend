package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// @Summary     Get githubProject
// @Description Get data from github_projects
// @ID          get-githubProject
// @Tags        githubProjects
// @Accept      json
// @Produce     json
// @Param       request body     request.GetByProjectUrl true "Project Url"
// @Success     200     {object} entity.GithubProject
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /edit/githubproject [get]
func (r *V1) getGithubProject(ctx *fiber.Ctx) error {

	projectUrl := ctx.Query("url")

	project, err := r.g.GetByProjectUrl(ctx.UserContext(), projectUrl)
	if err != nil {
		r.l.Error(err, "REST API - V1 - getGithubProject")

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(project)
}
