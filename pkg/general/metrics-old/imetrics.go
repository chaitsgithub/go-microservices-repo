package oldmetrics

// Counter is an interface for a metric that only increases.
type Counter interface {
	Inc()
	Add(float64)
}

// Gauge is an interface for a metric that can go up or down.
type Gauge interface {
	Inc()
	Dec()
	Add(float64)
	Sub(float64)
	Set(float64)
}

// Histogram is an interface for a metric that samples observations and puts them into buckets.
type Histogram interface {
	Observe(float64)
}

// Registry is the core interface that provides access to all metric types.
// A microservice should depend on this interface.
type Registry interface {
	Counter(name, help string, labels ...string) Counter
	Gauge(name, help string, labels ...string) Gauge
	Histogram(name, help string, labels ...string) Histogram
}
