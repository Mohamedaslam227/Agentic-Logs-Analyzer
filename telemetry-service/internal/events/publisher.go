package events

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"telemetry-service/internal/config"
	"telemetry-service/internal/detectors"

	"github.com/google/uuid"
)

type Publisher struct {
	client  *http.Client
	sinkURL string
	source  string
}

func NewPublisher(cfg *config.Config) *Publisher {
	return &Publisher{
		client: &http.Client{
			Timeout: cfg.EventTimeout,
		},
		sinkURL: cfg.EventSinkURL,
		source:  cfg.ServiceName,
	}
}

func (p *Publisher) Publish(ctx context.Context, signal *detectors.IncidentSignal) error {
	event := Event{
		ID:        uuid.NewString(),
		Type:      string(signal.Type),
		Severity:  string(signal.Severity),
		Namespace: signal.Namespace,
		Resource:  signal.Resource,
		Message:   signal.Message,
		Timestamp: signal.Timestamp,
		Metadata:  signal.Metadata,
		Source:    p.source,
	}
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.sinkURL, bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read and print the Agent's response
	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err == nil {
		fmt.Printf("\n--- Agent Response ---\n")
		fmt.Printf("Decision: %v\n", responseMap["decision"])
		fmt.Printf("Message: %v\n", responseMap["message"])
		fmt.Println("----------------------")
	} else {
		// Fallback if not JSON
		fmt.Println("Agent response received (non-JSON).")
	}

	return nil

}
