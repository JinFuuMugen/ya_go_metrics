package audit

import (
	"reflect"

	"github.com/JinFuuMugen/ya_go_metrics/internal/logger"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
)

// Publisher implements subscribe pattern, allowing multiple Observers to receive same events.

//generate:reset
//go:generate go run ../../cmd/reset/main.go
type Publisher struct {
	observers []Observer
}

// NewPublisher creates new Publisher instance without subscribers.
func NewPublisher() *Publisher {
	return &Publisher{}
}

// Subscribe registers an Observer to receive events.
func (p *Publisher) Subscribe(o Observer) {
	p.observers = append(p.observers, o)
}

// Publish sends given AuditEvent to subscribers.
func (p *Publisher) Publish(auditEvent models.AuditEvent) {
	for _, o := range p.observers {
		err := o.Notify(auditEvent)
		if err != nil {
			logger.Errorf("cannot send audit event %s to %s : %w", auditEvent, reflect.TypeOf(o), err)
		}
	}
}
