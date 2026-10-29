package project

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/Taya-san/portofolio-backend/internal/repo/project"

type tracedRepo struct {
	next ProjectRepo
}

func newTraced(next ProjectRepo) ProjectRepo {
	return &tracedRepo{next: next}
}

func startSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return otel.Tracer(_tracerName).Start(ctx, name, trace.WithAttributes(attrs...))
}

func endSpan(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	span.End()
}

func (pr *tracedRepo) Store(ctx context.Context, project *entity.Project) error {
	ctx, span := startSpan(ctx, "ProjectRepo.Store", attribute.String("project.url", project.ProjectUrl))

	err := pr.next.Store(ctx, project)
	endSpan(span, err)

	return err
}

func (pr *tracedRepo) GetByID(ctx context.Context, projectID string) (entity.Project, error) {
	ctx, span := startSpan(ctx, "ProjectRepo.GetByID", attribute.String("project.id", projectID))

	result, err := pr.next.GetByID(ctx, projectID)
	endSpan(span, err)

	return result, err
}

func (pr *tracedRepo) List(ctx context.Context, filter ProjectFilter) ([]entity.Project, int, error) {
	ctx, span := startSpan(ctx, "ProjectRepo.List", attribute.Bool("all", filter.All))

	results, n, err := pr.next.List(ctx, filter)
	endSpan(span, err)

	return results, n, err
}

func (pr *tracedRepo) Update(ctx context.Context, project *entity.Project) error {
	ctx, span := startSpan(ctx, "ProjectRepo.Update", attribute.String("project.id", project.ID))

	err := pr.next.Update(ctx, project)
	endSpan(span, err)

	return err
}

func (pr *tracedRepo) Delete(ctx context.Context, projectID string) error {
	ctx, span := startSpan(ctx, "ProjectRepo.Delete", attribute.String("project.id", projectID))

	err := pr.next.Delete(ctx, projectID)
	endSpan(span, err)

	return err
}
