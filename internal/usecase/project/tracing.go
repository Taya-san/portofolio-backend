package project

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/Taya-san/portofolio-backend/internal/usecase/project"

type tracedUseCase struct {
	next Project
}

func newTraced(next Project) Project {
	return &tracedUseCase{next: next}
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

func (p *tracedUseCase) Create(ctx context.Context, projectItem ProjectItem) (entity.Project, error) {
	ctx, span := startSpan(ctx, "ProjectUseCase.Create", attribute.String("project.url", projectItem.ProjectUrl))

	result, err := p.next.Create(ctx, projectItem)
	endSpan(span, err)

	return result, err
}

func (p *tracedUseCase) Get(ctx context.Context, projectID string) (entity.Project, error) {
	ctx, span := startSpan(ctx, "ProjectUseCase.Get", attribute.String("project.id", projectID))

	result, err := p.next.Get(ctx, projectID)
	endSpan(span, err)

	return result, err
}

func (p *tracedUseCase) List(ctx context.Context, all bool) ([]entity.Project, int, error) {
	ctx, span := startSpan(ctx, "ProjectUseCase.List", attribute.Bool("all", all))

	results, n, err := p.next.List(ctx, all)
	endSpan(span, err)

	return results, n, err
}

func (p *tracedUseCase) Update(ctx context.Context, projectID string, projectItem ProjectItem) (entity.Project, error) {
	ctx, span := startSpan(ctx, "ProjectUseCase.Update", attribute.String("project.id", projectID))

	result, err := p.next.Update(ctx, projectID, projectItem)
	endSpan(span, err)

	return result, err
}

func (p *tracedUseCase) Delete(ctx context.Context, projectID string) error {
	ctx, span := startSpan(ctx, "ProjectUseCase.Delete", attribute.String("project.id", projectID))

	err := p.next.Delete(ctx, projectID)
	endSpan(span, err)

	return err
}
