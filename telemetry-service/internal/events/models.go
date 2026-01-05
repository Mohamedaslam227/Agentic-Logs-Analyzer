package events
import "time"

type Event struct {
	ID string `json:"id,omitempty"`
	Type string `json:"type"`
	Severity string `json:"severity"`
	Namespace string `json:"namespace,omitempty"`
	Resource string `json:"resource"`
	Message string `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Source string `json:"source"`
}