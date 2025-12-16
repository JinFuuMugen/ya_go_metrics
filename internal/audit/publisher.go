package audit

import "github.com/JinFuuMugen/ya_go_metrics/internal/models"

type Publisher struct {
	observers []Observer
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Subscribe(o Observer) {
	p.observers = append(p.observers, o)
}

func (p *Publisher) Publish(auditEvent models.AuditEvent) {
	for _, o := range p.observers {
		_ = o.Notify(auditEvent)
	}
}
