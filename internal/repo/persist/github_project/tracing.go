package githubproject

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/Taya-san/portofolio-backend/internal/repo/persist/github_project"

type tracedRepo struct {
	next GithubProjectRepo
}

func newTraced(next GithubProjectRepo) GithubProjectRepo {
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

func (pr *tracedRepo) UpStore(ctx context.Context, project *entity.GithubProject) error {
	ctx, span := startSpan(ctx, "GithubProjectRepo.Store", attribute.String("project.url", project.ProjectUrl))

	err := pr.next.UpStore(ctx, project)
	endSpan(span, err)

	return err
}

func (pr *tracedRepo) GetByProjectUrl(ctx context.Context, projectUrl string) (entity.GithubProject, error) {
	ctx, span := startSpan(ctx, "GithubProjectRepo.GetByID", attribute.String("project.url", projectUrl))

	result, err := pr.next.GetByProjectUrl(ctx, projectUrl)
	endSpan(span, err)

	return result, err
}

func (pr *tracedRepo) List(ctx context.Context) ([]entity.GithubProject, int, error) {
	ctx, span := startSpan(ctx, "GithubProjectRepo.List")

	results, n, err := pr.next.List(ctx)
	endSpan(span, err)

	return results, n, err
}

func (pr *tracedRepo) Delete(ctx context.Context, projectID string) error {
	ctx, span := startSpan(ctx, "GithubProjectRepo.Delete", attribute.String("project.id", projectID))

	err := pr.next.Delete(ctx, projectID)
	endSpan(span, err)

	return err
}
