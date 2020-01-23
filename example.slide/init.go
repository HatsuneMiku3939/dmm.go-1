import (
	"contrib.go.opencensus.io/exporter/jaeger"
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
}
