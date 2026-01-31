package audit

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
)

// FileObserver implements audit event observation by writing to a file.
type FileObserver struct {
	file *os.File
}

// NewFileObserver creates a new FileObserver that writes to a file with a provided path.
// If file does not exist - it will be created.
func NewFileObserver(path string) (*FileObserver, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}

	return &FileObserver{file: f}, nil
}

// Notify writes the given audit to a file.
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
