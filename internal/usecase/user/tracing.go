package user

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const _tracerName = "github.com/Taya-san/portofolio-backend/internal/usecase/user"

type tracedUseCase struct {
	next User
}

func newTraced(next User) User {
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

func (u *tracedUseCase) Login(ctx context.Context, userItem UserItem) (string, error) {
	ctx, span := startSpan(ctx, "UserUseCase.Login")

	result, err := u.next.Login(ctx, userItem)
	endSpan(span, err)

	return result, err
}
