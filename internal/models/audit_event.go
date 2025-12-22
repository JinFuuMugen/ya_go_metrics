package models

// AuditEvent represents a single audit record for metric update operations.
type AuditEvent struct {
	// TS is a Unix timestamp when the audit event occurred.
	TS int64 `json:"ts"`

	// Metrics contains metrics affected by the operation.
	Metrics []string `json:"metrics"`

	// IPAddress is the IP address of the client that triggered the event.
	IPAddress string `json:"ip_address"`
}
