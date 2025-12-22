package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
)

// HTTPObserver implements audit event observation by sending events to a remote HTTP endpoint.
type HTTPObserver struct {
	url    string
	client *http.Client
}

// NewHTTPObserver creates a new HTTPObserver that sends audit events to given URL.
func NewHTTPObserver(url string) *HTTPObserver {
	return &HTTPObserver{
		url: url,
		client: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// Notify sends the given AuditEvent via POST request with Content-Type:application/json header and JSON-encoded body
func (ho *HTTPObserver) Notify(auditEvent models.AuditEvent) error {
	data, err := json.Marshal(auditEvent)
	if err != nil {
		return fmt.Errorf("cannot marshal audit event: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, ho.url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := ho.client.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}
