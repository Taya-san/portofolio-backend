package user

import (
	"context"

	"github.com/Taya-san/portofolio-backend/internal/entity"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/Taya-san/portofolio-backend/internal/repo/user"

type tracedRepo struct {
	next UserRepo
}

func newTraced(next UserRepo) UserRepo {
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

func (ur *tracedRepo) ListPermittedPasswordHashes(ctx context.Context) ([]entity.User, error) {
	ctx, span := startSpan(ctx, "UserRepo.ListPermittedPasswordHashes")

	results, err := ur.next.ListPermittedPasswordHashes(ctx)
	endSpan(span, err)

	return results, err
}
