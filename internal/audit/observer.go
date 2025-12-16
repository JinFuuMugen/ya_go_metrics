package audit

import "github.com/JinFuuMugen/ya_go_metrics/internal/models"

type Observer interface {
	Notify(auditEvent models.AuditEvent) error
}
