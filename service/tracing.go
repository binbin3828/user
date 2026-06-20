package service

import (
	"context"
	"os"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gin-gonic/gin"
)

func newTraceExporter() (sdktrace.SpanExporter, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint != "" {
		return otlptracegrpc.New(context.Background(),
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(),
		)
	}
	return stdouttrace.New(stdouttrace.WithWriter(os.Stdout))
}

func sampleRate() float64 {
	const defaultRate = 0.1
	s := os.Getenv("OTEL_TRACE_SAMPLE_RATE")
	if s == "" {
		return defaultRate
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return defaultRate
	}
	if f < 0 || f > 1 {
		return defaultRate
	}
	return f
}

func InitTracerProvider(serviceName string) (*sdktrace.TracerProvider, error) {
	exp, err := newTraceExporter()
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp, sdktrace.WithBatchTimeout(time.Second)),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(sampleRate())),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func TracingMiddleware(serviceName string) gin.HandlerFunc {
	return otelgin.Middleware(serviceName)
}

func traceIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return spanCtx.TraceID().String()
	}
	return ""
}

func spanIDFromContext(ctx context.Context) string {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasSpanID() {
		return spanCtx.SpanID().String()
	}
	return ""
}
