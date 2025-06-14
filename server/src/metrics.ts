import { getMeter } from "./observability.js"

const meter = getMeter()

// Create a Counter
export const httpRequestsCounter = meter.createCounter("http.requests", {
    description: "Count of HTTP requests processed",
    unit: "1",
})

// Create a Histogram
export const httpResponseTimeHistogram = meter.createHistogram("http.response.duration", {
    description: "Duration of HTTP requests",
    unit: "ms",
})

// Create an Up/Down Counter
export const httpActiveConnectionsCounter = meter.createUpDownCounter("http.active_connections", {
    description: "Number of active HTTP connections",
    unit: "1",
})

// Create a Gauge using Observable Counter
meter
    .createObservableGauge("system.memory.usage", { description: "Current memory usage", unit: "bytes" })
    .addCallback((callback) => {
        const memoryUsage = process.memoryUsage()
        callback.observe(memoryUsage.heapUsed, { type: "heap.used" })
        callback.observe(memoryUsage.heapTotal, { type: "heap.total" })
        callback.observe(memoryUsage.rss, { type: "rss" })
    })
