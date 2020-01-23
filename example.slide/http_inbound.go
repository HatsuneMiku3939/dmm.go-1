package main

import (
	"github.com/HatsuneMiku3939/ocecho"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/trace"
)

func main() {
	if err := InitJaegerTrace("http-service", 1.0); err != nil {
		panic(err)
	}

	// Echo instance
	e := echo.New()
	e.Use(ocecho.OpenCensusMiddleware(
		ocecho.OpenCensusConfig{
			TraceOptions: ocecho.TraceOptions{IsPublicEndpoint: true},
		},
	))

}
