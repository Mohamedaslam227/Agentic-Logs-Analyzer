package detectors
import "time"

type Severity string
const (
	SeverityLow Severity = "low"
	SeverityMedium Severity = "medium"
	SeverityHigh Severity = "high"
	SeverityCritical Severity = "critical"
)

type SignalType string

const (
	SignalCrashLoop SignalType = "crash_loop"
	SignalOOM SignalType = "oom"
	SignalCPUSpike SignalType = "cpu_spike"
	SignalAnamoly SignalType = "anamoly"
)


type IncidentSignal struct {
	ID string `json:"id,omitempty"`
	Type SignalType `json:"type"`
	Severity Severity `json:"severity"`
	Namespace string `json:"namespace,omitempty"`
	Resource  string `json:"resource"`
	Message string `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Metadata map[string]string `json:"metadata,omitempty"`

}

type SignalInput struct {
	Metrics map[string][]float64
	Labels map[string]string
}

type Detector interface {
	Name() string
	Detect(input SignalInput) (*IncidentSignal, bool)
}