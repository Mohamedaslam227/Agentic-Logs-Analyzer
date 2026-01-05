package detectors
import (
	"fmt"
	"strings"
	"time"
)

type CPUSpikeDetector struct {
	Threshold float64
}

func NewCPUSpikeDetector(threshold float64) *CPUSpikeDetector {
	return &CPUSpikeDetector{
		Threshold: threshold,
	}
}

func (d *CPUSpikeDetector) Name() string {
	return "cpu_spike_detector"
}


func (d *CPUSpikeDetector) Detect(input SignalInput) (*IncidentSignal, bool) {
	for key, values := range input.Metrics {
		if !strings.HasPrefix(key,"cpu:") {
			continue
		}
		avg := average(values)
		if avg > d.Threshold {
			namespace, resource := parseKey(key)
			signal := &IncidentSignal{
				Type: SignalCPUSpike,
				Severity: ClassifySeverity(avg, d.Threshold),
				Namespace: namespace,
				Resource: resource,
				Message: fmt.Sprintf(
					"CPU spike detected: average usage %.2f millicores exceeds threshold %.2f",
					avg,
					d.Threshold,
				),
				Timestamp: time.Now(),
				Metadata: map[string]string{
					"average_cpu_millicores": fmt.Sprintf("%.2f", avg),
					"threshold_millicores": fmt.Sprintf("%.2f", d.Threshold),
				},
			}
			return signal, true
		}
	}
	return nil, false
}


func average(values []float64) float64 {
	if len(values) == 0{
		return 0
	}
	var sum float64
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}

func parseKey(key string) (string, string) {
	parts := strings.Split(key, ":")

	if len(parts) == 3 {
		return parts[1], parts[2]

	}
	if len(parts) == 2 {
		return "", parts[1]
	}

	return "", "unknown"
}

func ClassifySeverity(avg, threshold float64) Severity {
	switch {
	case avg >= threshold*1.5:
		return SeverityCritical
	case avg >= threshold*1.2:
		return SeverityHigh
	case avg >= threshold*1.1:
		return SeverityMedium
	default:
		return SeverityLow
	}
}