package opentelemetry

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
)

func TestOpenTelemetry(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer("otel")
	_, span := tracer.Start(context.Background(), "zx")
	defer span.End()
	
}

func initZipkin(t *testing.T) {
	exporter, err := zipkin.New("http://124.222.29.179:19411/api/v2/spans")
	if err != nil {
		t.Fatal(err)
	}
	batcher := sdktrace.NewBatchSpanProcessor(exporter)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("opentelemetry-demo"),
		)),
	)
	otel.SetTracerProvider(tp)
}