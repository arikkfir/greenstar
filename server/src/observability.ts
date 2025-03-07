import { NodeSDK } from "@opentelemetry/sdk-node"
import { detectResources, envDetector } from "@opentelemetry/resources"
import { BatchSpanProcessor } from "@opentelemetry/sdk-trace-node"
import { OTLPTraceExporter } from "@opentelemetry/exporter-trace-otlp-http"
import { PrometheusExporter } from "@opentelemetry/exporter-prometheus"
import { getNodeAutoInstrumentations } from "@opentelemetry/auto-instrumentations-node"
import { diag, DiagConsoleLogger, DiagLogLevel, metrics, trace } from "@opentelemetry/api"
import { RuntimeNodeInstrumentation } from "@opentelemetry/instrumentation-runtime-node"
import { NetInstrumentation } from "@opentelemetry/instrumentation-net"
import { HttpInstrumentation } from "@opentelemetry/instrumentation-http"
import { FsInstrumentation } from "@opentelemetry/instrumentation-fs"
import { ExpressInstrumentation } from "@opentelemetry/instrumentation-express"
import { DnsInstrumentation } from "@opentelemetry/instrumentation-dns"

// Validate required configuration is provided
if (!process.env.OTEL_EXPORTER_PROMETHEUS_PORT) {
    throw new Error("environment variable OTEL_EXPORTER_PROMETHEUS_PORT is required")
}
if (!process.env.OTEL_EXPORTER_PROMETHEUS_ENDPOINT) {
    throw new Error("environment variable OTEL_EXPORTER_PROMETHEUS_ENDPOINT is required")
}
if (!process.env.OTEL_EXPORTER_OTLP_ENDPOINT) {
    throw new Error("environment variable OTEL_EXPORTER_OTLP_ENDPOINT is required")
}

// Enable debug logging (optional)
diag.setLogger(new DiagConsoleLogger(), DiagLogLevel.INFO)

// Automatically detect resource from environment variables
const resource = detectResources({ detectors: [envDetector] })

// Configure Prometheus exporter
const prometheusExporter = new PrometheusExporter({
    port: parseInt(process.env.OTEL_EXPORTER_PROMETHEUS_PORT),
    endpoint: process.env.OTEL_EXPORTER_PROMETHEUS_ENDPOINT,
})

// Configure OTLP trace exporter
const traceExporter = new OTLPTraceExporter({ url: `${process.env.OTEL_EXPORTER_OTLP_ENDPOINT}/v1/traces` })
const spanProcessor = new BatchSpanProcessor(traceExporter)

// Create and initialize the OpenTelemetry SDK with the updated spanProcessors array
const sdk = new NodeSDK({
    resource,
    spanProcessors: [spanProcessor],
    metricReader: prometheusExporter,
    instrumentations: [
        new DnsInstrumentation(),
        new ExpressInstrumentation(),
        new FsInstrumentation(),
        new HttpInstrumentation(),
        new NetInstrumentation(),
        new RuntimeNodeInstrumentation(),
    ],
})

// Initialize the SDK, and ensure it shuts down when the process exits
sdk.start()
process.on("SIGTERM", async () => {
    try {
        await sdk.shutdown()
        console.log("OpenTelemetry SDK shut down")
    } catch (e) {
        console.error("Error shutting down OpenTelemetry SDK", e)
    }
})

// Provide the default meter to the rest of the app
// A "meter" here is the equivalent of "logger", i.e. you use one "meter" to create many metrics
export function getMeter() {
    return metrics.getMeter("default-meter")
}

// Provide the default trace to the rest of the app
// A "trace" here is the equivalent of "logger", i.e. you use one "trace" to create many traces
export function getTracer() {
    return trace.getTracer("default-tracer")
}
