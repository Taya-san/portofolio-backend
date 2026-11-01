package githubproject

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/Taya-san/portofolio-backend/internal/usecase/github_project"

type tracedUseCase struct {
	next GithubProject
}

func newTraced(next GithubProject) GithubProject {
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

func (p *tracedUseCase) RefreshMirror(ctx context.Context) error {
	ctx, span := startSpan(ctx, "GithubProjectUseCase.RefreshMirror")

	err := p.next.RefreshMirror(ctx)
	endSpan(span, err)

	return err
}
