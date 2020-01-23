package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"contrib.go.opencensus.io/exporter/jaeger"
	"contrib.go.opencensus.io/exporter/prometheus"
	"github.com/HatsuneMiku3939/ocecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

// InitJaeger Initialize opencensus Jaeger trace exporter
func InitJaegerTrace(serviceName string, sampleRate float64) error {
	// Port details: https://www.jaegertracing.io/docs/getting-started/
	agentEndpointURI := os.Getenv("JAEGER_AGENT_ENDPOINT")
	collectorEndpointURI := os.Getenv("JAEGER_COLLECTOR_ENDPOINT")

	je, err := jaeger.NewExporter(jaeger.Options{
		AgentEndpoint:     agentEndpointURI,
		CollectorEndpoint: collectorEndpointURI,
		ServiceName:       serviceName,
	})
	if err != nil {
		return fmt.Errorf("Failed to create the Jaeger exporter: %v", err)
	}

	// Set sampleRate
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(sampleRate)})

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(je)
	return nil
}

// InitPrometheus Initialize opencensus Prometheus metric exporter
func InitPrometheus() error {
	pe, err := prometheus.NewExporter(prometheus.Options{})
	if err != nil {
		return fmt.Errorf("Failed to create the Prometheus exporter: %v", err)
	}
	if err := view.Register(ocecho.DefaultServerViews...); err != nil {
		return fmt.Errorf("Failed to register server metric view: %v", err)
	}

	// Register stats and trace exporters to export
	// the collected data.
	view.RegisterExporter(pe)

	// Start Prometheus endpoint
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", pe)
		http.ListenAndServe(":8888", mux)
	}()
	return nil
}

// TraceRequest Create http.Client, http.Request with trace context
func TraceRequest(ctx context.Context, method, url string, body io.Reader) (*http.Client, *http.Request, error) {
	// trace client
	client := &http.Client{Transport: &ochttp.Transport{}}

	// wrap request with context
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	req = req.WithContext(ctx)

	return client, req, nil
}

func main() {
	if err := InitJaegerTrace("service_a", 1.0); err != nil {
		panic(err)
	}
	if err := InitPrometheus(); err != nil {
		panic(err)
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(ocecho.OpenCensusMiddleware(
		ocecho.OpenCensusConfig{
			TraceOptions: ocecho.TraceOptions{IsPublicEndpoint: true},
		},
	))

	// GET /api/trace
	e.GET("/api/trace", func(c echo.Context) error {
		ctx := c.Request().Context()

		client, req, err := TraceRequest(ctx, "GET", "http://service_b:8080/api/internal", nil)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		defer resp.Body.Close()

		return c.Stream(http.StatusOK, "text/plain", resp.Body)
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
