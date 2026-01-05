package metrics
import (
	"fmt"
	"telemetry-service/internal/detectors"
)

func AggregateMetrics(metrics []Metric) detectors.SignalInput {
	aggregated := make(map[string][]float64)
	labels := make(map[string]string)

	for _, m := range metrics {
		key := buildMetricKey(m)
		aggregated[key] = append(aggregated[key], m.Value)
		for k,v := range m.Labels {
			labels[k] = v
		}
	}
	return detectors.SignalInput{
		Metrics: aggregated,
		Labels: labels,
	}

	
}

func buildMetricKey(m Metric) string {
	if m.Namespace !="" {
		return fmt.Sprintf("%s:%s:%s", m.Type, m.Namespace, m.Resource)
	}
	return fmt.Sprintf("%s:%s", m.Type, m.Resource)
}