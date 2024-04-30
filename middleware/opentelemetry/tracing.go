package opentelemetry

import (
	"zweb"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

func (m *MiddlewareBuilder) Build() zweb.Middleware {
	return func(next zweb.HandleFunc) zweb.HandleFunc {
		return func(ctx *zweb.Context) {
			reqCtx := ctx.Req.Context()
			reqCtx, span := m.Tracer.Start(reqCtx, "unknow", trace.WithAttributes())
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			ctx.Req = ctx.Req.WithContext(reqCtx)
			next(ctx)
			if ctx.MatchedRoute != "" {
				span.SetName(string(ctx.MatchedRoute))
			}
			span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
		}
	}
}