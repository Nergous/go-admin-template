package observability

import (
	"context"
	"log/slog"
	"time"
)

type contextKey string

const requestMetadataContextKey contextKey = "request_metadata"

type RequestMetadata struct {
	RequestID string
	StartedAt time.Time
	ClientIP  string
	UserAgent string
	Device    string
}

func WithRequestMetadata(ctx context.Context, metadata RequestMetadata) context.Context {
	if metadata.StartedAt.IsZero() {
		metadata.StartedAt = time.Now()
	}

	return context.WithValue(ctx, requestMetadataContextKey, metadata)
}

func RequestMetadataFromContext(ctx context.Context) RequestMetadata {
	metadata, _ := ctx.Value(requestMetadataContextKey).(RequestMetadata)
	return metadata
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	metadata := RequestMetadataFromContext(ctx)
	metadata.RequestID = requestID
	return WithRequestMetadata(ctx, metadata)
}

func RequestIDFromContext(ctx context.Context) string {
	return RequestMetadataFromContext(ctx).RequestID
}

type contextHandler struct {
	slog.Handler
}

func (h contextHandler) Handle(ctx context.Context, record slog.Record) error {
	metadata := RequestMetadataFromContext(ctx)
	if metadata.RequestID != "" {
		record.AddAttrs(slog.String("request_id", metadata.RequestID))
	}
	if !metadata.StartedAt.IsZero() {
		record.AddAttrs(
			slog.Time("request_started_at", metadata.StartedAt),
			slog.Duration("request_duration", time.Since(metadata.StartedAt)),
		)
	}
	if metadata.ClientIP != "" {
		record.AddAttrs(slog.String("client_ip", metadata.ClientIP))
	}
	if metadata.UserAgent != "" {
		record.AddAttrs(slog.String("user_agent", metadata.UserAgent))
	}
	if metadata.Device != "" {
		record.AddAttrs(slog.String("device", metadata.Device))
	}

	return h.Handler.Handle(ctx, record)
}

func (h contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return contextHandler{Handler: h.Handler.WithAttrs(attrs)}
}

func (h contextHandler) WithGroup(name string) slog.Handler {
	return contextHandler{Handler: h.Handler.WithGroup(name)}
}
