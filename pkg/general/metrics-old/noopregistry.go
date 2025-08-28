package oldmetrics

type noopRegistry struct{}

func NewNoopRegistry() *noopRegistry {
	return &noopRegistry{}
}

type noopCounter struct{}

func (c *noopCounter) Inc()        {}
func (c *noopCounter) Add(float64) {}

type noopGauge struct{}

func (g *noopGauge) Inc()        {}
func (g *noopGauge) Dec()        {}
func (g *noopGauge) Add(float64) {}
func (g *noopGauge) Sub(float64) {}
func (g *noopGauge) Set(float64) {}

type noopHistogram struct{}

func (h *noopHistogram) Observe(float64) {}

func (n *noopRegistry) Counter(name, help string, labels ...string) Counter {
	return &noopCounter{}
}
func (n *noopRegistry) Gauge(name, help string, labels ...string) Gauge {
	return &noopGauge{}
}
func (n *noopRegistry) Histogram(name, help string, labels ...string) Histogram {
	return &noopHistogram{}
}
