package audit

import "github.com/JinFuuMugen/ya_go_metrics/internal/models"

// Observer defines the interface for receiving AuditEvents.
type Observer interface {
	//Notify sends single AuditEvent to observer.
	Notify(auditEvent models.AuditEvent) error
}
