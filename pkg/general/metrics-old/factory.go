package oldmetrics

import "fmt"

const (
	METRICS_PROMETHEUS = "prometheus"
	METRICS_NOOP       = "noop"
)

func NewMetricsRegistry(backendType string) (Registry, error) {
	switch backendType {
	case METRICS_PROMETHEUS:
		return NewPrometheusRegistry(), nil
	case METRICS_NOOP:
		return NewNoopRegistry(), nil
	default:
		return nil, fmt.Errorf("not a supported metrics registry")
	}
}
