package opentelemetry

import (
	"zweb"
	"go.opentelemetry.io/otel"
	"testing"
	"time"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer("")
	initZipkin(t)
	s := zweb.NewHTTPServer()
	s.Get("/", func(ctx *zweb.Context) {
		// ctx.Resp.Write([]byte("hello, world2"))
		ctx.Resp.Write([]byte("hello, world2"))
	})
	s.Get("/user", func(ctx *zweb.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()
		c, third2 := tracer.Start(c, "third_layer_1")
		time.Sleep(300 * time.Millisecond)
		third2.End()
		second.End()
		// ctx.RespStatusCode = 200
		// ctx.RespData = []byte("hello, world")
	})

	s.Use((&MiddlewareBuilder{Tracer: tracer}).Build())
	s.Start(":8081")
}
