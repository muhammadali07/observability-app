package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	tracer       trace.Tracer
	otlpEndpoint string
)

func init() {
	otlpEndpoint = os.Getenv("OTLP_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = "localhost:4318"
	}
}

// newConsoleExporter creates a new exporter for testing, outputting to the console.
func newConsoleExporter() (sdktrace.SpanExporter, error) {
	return stdouttrace.New()
}

// newOTLPExporter creates a new OTLP exporter with the specified context.
func newOTLPExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	insecureOpt := otlptracehttp.WithInsecure()
	endpointOpt := otlptracehttp.WithEndpoint(otlpEndpoint)
	return otlptracehttp.New(ctx, insecureOpt, endpointOpt)
}

// newTraceProvider creates a new TracerProvider with the specified exporter.
func newTraceProvider(exp sdktrace.SpanExporter) *sdktrace.TracerProvider {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("myapp"),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
}

func main() {
	ctx := context.Background()
	exp, err := newOTLPExporter(ctx)
	if err != nil {
		log.Fatalf("failed to initialize exporter: %v", err)
	}

	tp := newTraceProvider(exp)
	defer func() { _ = tp.Shutdown(ctx) }()

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("myapp")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/devices", getDevices)

	log.Fatal(http.ListenAndServe(":8080", r))
}

func getDevices(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "HTTP GET /devices")
	defer span.End()

	db(ctx)
	time.Sleep(1 * time.Second)
	w.Write([]byte("ok"))
}

func db(ctx context.Context) {
	_, span := tracer.Start(ctx, "SQL SELECT")
	defer span.End()
	time.Sleep(2 * time.Second)
}
