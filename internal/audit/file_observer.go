package audit

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
)

type FileObserver struct {
	file *os.File
}

func NewFileObserver(path string) (*FileObserver, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return &FileObserver{file: f}, nil
}

func (fo *FileObserver) Notify(auditEvent models.AuditEvent) error {
	data, err := json.Marshal(auditEvent)
	if err != nil {
		return fmt.Errorf("cannot marshal audit event")
	}

	_, err = fo.file.Write(append(data, '\n'))
	if err != nil {
		return fmt.Errorf("cannot write event to file: %w", err)
	}

	return nil
}
