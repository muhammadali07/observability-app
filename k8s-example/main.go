package main

import (
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo-contrib/jaegertracing"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

// urlSkipper ignores metrics route on some middleware
func urlSkipper(c echo.Context) bool {
	if strings.HasPrefix(c.Path(), "/testurl") {
		return true
	}
	return false
}

func main() {
	e := echo.New()
	c := jaegertracing.New(e, urlSkipper)
	defer c.Close()

	// Routing
	e.GET("/", helloHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

// Handler untuk rute "/"
func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from your Echo backend!")
}

// Inisialisasi Jaeger Tracer
func initJaeger() (opentracing.Tracer, io.Closer) {
	cfg := config.Configuration{
		ServiceName: "nama-service-anda",
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: true,
		},
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		log.Fatal("Failed to initialize Jaeger tracer:", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}
