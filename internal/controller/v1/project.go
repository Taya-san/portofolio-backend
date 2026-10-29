package v1

import (
	"errors"
	"net/http"

	"github.com/Taya-san/portofolio-backend/internal/controller/v1/request"
	"github.com/Taya-san/portofolio-backend/internal/controller/v1/response"
	"github.com/Taya-san/portofolio-backend/internal/entity"
	"github.com/Taya-san/portofolio-backend/internal/usecase/project"
	"github.com/gofiber/fiber/v2"
)

// @Summary     Create project
// @Description Create a new project
// @ID          create-project
// @Tags        projects
// @Accept      json
// @Produce     json
// @Param       request body     request.CreateProject true "Project data"
// @Success     201     {object} entity.Project
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /edit [post]
func (r *V1) createProject(ctx *fiber.Ctx) error {
	var body request.CreateProject

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "REST API - V1 - createProject")

		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "REST API - V1 - createProject")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	project, err := r.p.Create(ctx.UserContext(),
		project.ProjectItem{
			ProjectUrl: body.ProjectUrl,
			Title:      body.Title,
			Languages:  body.Languages,
			Visible:    body.Visible,
			Custom:     body.Custom,
			Body:       body.Body,
		},
	)
	if err != nil {
		r.l.Error(err, "REST API - V1 - createProject")

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusCreated).JSON(project)
}

// @Summary     List projects
// @Description List projects with filtering
// @ID          list-projects
// @Tags        projects
// @Produce     json
// @Success     200    {object} response.ProjectList
// @Failure     401    {object} response.Error
// @Failure     500    {object} response.Error
// @Security    BearerAuth
// @Router      /projects [get]
func (r *V1) listProjects(ctx *fiber.Ctx) error {
	_, ok := ctx.Locals("userID").(string)
	all := true
	if !ok {
		all = false
	}

	projects, n, err := r.p.List(ctx.UserContext(), all)
	if err != nil {
		r.l.Error(err, "REST API - V1 - listProjects")

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(response.ProjectList{
		Projects: projects,
		Total:    n,
	})
}

// @Summary     Get project
// @Description Get a project by ID
// @ID          get-project
// @Tags        projects
// @Produce     json
// @Param       id  path     string true "Project ID"
// @Success     200 {object} entity.Project
// @Failure     401 {object} response.Error
// @Failure     403 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /project/{id} [get]
func (r *V1) getProject(ctx *fiber.Ctx) error {
	_, ok := ctx.Locals("userID").(string)
	projectID := ctx.Params("id")
	project, err := r.p.Get(ctx.UserContext(), projectID)
	if err != nil {
		r.l.Error(err, "REST API - V1 - getProject")

		if errors.Is(err, entity.ErrProjectNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "project not found")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}
	if !ok && !project.Visible {
		return errorResponse(ctx, http.StatusForbidden, "forbidden")
	}

	return ctx.Status(http.StatusOK).JSON(project)
}

// @Summary     Update project
// @Description Update project data
// @ID          update-project
// @Tags        projects
// @Accept      json
// @Produce     json
// @Param       id      path     string            true "Project ID"
// @Param       request body     request.UpdateProject  true "Updated project data"
// @Success     200     {object} entity.Project
// @Failure     400     {object} response.Error
// @Failure     401     {object} response.Error
// @Failure     403     {object} response.Error
// @Failure     404     {object} response.Error
// @Failure     500     {object} response.Error
// @Security    BearerAuth
// @Router      /edit/{id} [put]
func (r *V1) updateProject(ctx *fiber.Ctx) error {
	_, ok := ctx.Locals("userID").(string)
	if !ok {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	projectID := ctx.Params("id")

	var body request.UpdateProject

	if err := ctx.BodyParser(&body); err != nil {
		r.l.Error(err, "REST API - v1 - updateProject")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	if err := r.v.Struct(body); err != nil {
		r.l.Error(err, "REST API - V1 - updateProject")

		return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
	}

	project, err := r.p.Update(ctx.UserContext(), projectID,
		project.ProjectItem{
			ProjectUrl: body.ProjectUrl,
			Title:      body.Title,
			Languages:  body.Languages,
			Visible:    body.Visible,
			Custom:     body.Custom,
			Body:       body.Body,
		},
	)
	if err != nil {
		r.l.Error(err, "REST API - V1 - updateProject")

		if errors.Is(err, entity.ErrProjectNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "project not found")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.Status(http.StatusOK).JSON(project)
}

// @Summary     Delete project
// @Description Delete a project by ID
// @ID          delete-project
// @Tags        projects
// @Param       id  path     string true "Project ID"
// @Success     204 "No Content"
// @Failure     401 {object} response.Error
// @Failure     404 {object} response.Error
// @Failure     500 {object} response.Error
// @Security    BearerAuth
// @Router      /edit/{id} [delete]
func (r *V1) deleteProject(ctx *fiber.Ctx) error {
	_, ok := ctx.Locals("userID").(string)
	if !ok {
		return errorResponse(ctx, http.StatusUnauthorized, "unauthorized")
	}

	projectID := ctx.Params("id")

	err := r.p.Delete(ctx.UserContext(), projectID)
	if err != nil {
		r.l.Error(err, "restapi - v1 - deleteProject")

		if errors.Is(err, entity.ErrProjectNotFound) {
			return errorResponse(ctx, http.StatusNotFound, "project not found")
		}

		return errorResponse(ctx, http.StatusInternalServerError, "internal server error")
	}

	return ctx.SendStatus(http.StatusNoContent)
}
