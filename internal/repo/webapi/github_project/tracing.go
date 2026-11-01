package githubproject

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/Taya-san/portofolio-backend/internal/repo/webapi/github_project"

type tracedRepo struct {
	next GithubProjectApi
}

func newTraced(next GithubProjectApi) GithubProjectApi {
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

func (pr *tracedRepo) Fetch(ctx context.Context) ([]entity.GithubProject, error) {
	ctx, span := startSpan(ctx, "GithubProjectApi.Store")

	allUserProjectsData, err := pr.next.Fetch(ctx)
	endSpan(span, err)

	return allUserProjectsData, err
}
