package sender

import (
	"bytes"
	"compress/gzip"
	"crypto/rsa"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/JinFuuMugen/ya_go_metrics/internal/config"
	"github.com/JinFuuMugen/ya_go_metrics/internal/cryptography"
	"github.com/JinFuuMugen/ya_go_metrics/internal/cryptography/rsacrypto"
	"github.com/JinFuuMugen/ya_go_metrics/internal/models"
	"github.com/JinFuuMugen/ya_go_metrics/internal/network"
	"github.com/JinFuuMugen/ya_go_metrics/internal/pool"
	"github.com/JinFuuMugen/ya_go_metrics/internal/storage"
	"github.com/go-resty/resty/v2"
)

var bufferPool = pool.New(func() *bytes.Buffer {
	return &bytes.Buffer{}
})

// Sender defines an interface for sending collected metrics.
type Sender interface {
	Process([]storage.Counter, []storage.Gauge) error
	Compress(data []byte) ([]byte, error)
}

type values struct {
	addr      string
	client    *resty.Client
	key       string
	publicKey *rsa.PublicKey
}

// NewSender creates a new Sender instance using the provided configuration.
func NewSender(cfg config.AgentConfig, publicKey *rsa.PublicKey) *values {
	return &values{cfg.Addr, resty.New(), cfg.Key, publicKey}
}

// Compress compresses data using gzip algorithm.
func (v *values) Compress(data []byte) ([]byte, error) {

	b := bufferPool.Get()
	b.Reset()
	defer bufferPool.Put(b)

	w := gzip.NewWriter(b)

	if _, err := w.Write(data); err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}

	compressed := make([]byte, b.Len())
	copy(compressed, b.Bytes())
	return compressed, nil
}

// Process serializes metrics, compresses them and sends to the server.
func (v *values) Process(counters []storage.Counter, gauges []storage.Gauge) error {
	var err error

	var metrics []models.Metrics

	for _, c := range counters {
		cDelta := c.GetValue().(int64)
		metrics = append(metrics, models.Metrics{
			ID:    c.GetName(),
			MType: c.GetType(),
			Delta: &cDelta,
			Value: nil,
		})
	}
	for _, g := range gauges {
		gValue := g.GetValue().(float64)
		metrics = append(metrics, models.Metrics{
			ID:    g.GetName(),
			MType: g.GetType(),
			Delta: nil,
			Value: &gValue,
		})
	}
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("cannot serialize metric to json: %w", err)
	}
	compressedData, err := v.Compress(jsonData)
	if err != nil {
		return fmt.Errorf("error while compressing data: %w", err)
	}

	dataToSend := compressedData
	encrypted := false

	if v.publicKey != nil {
		dataToSend, err = rsacrypto.Encrypt(v.publicKey, compressedData)
		if err != nil {
			return fmt.Errorf("failed encrypt data: %w", err)
		}
		encrypted = true
	}
	url := "http://" + v.addr + "/updates/"

	req := v.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Content-Encoding", "gzip").
		SetBody(dataToSend)

	if encrypted {
		req.SetHeader("X-Encrypted", "rsa")
	}

	if v.key != "" {
		hash := cryptography.GetHMACSHA256(jsonData, v.key)
		req.SetHeader("HashSHA256", hex.EncodeToString(hash))
	}

	ip, err := network.OutboundIPTo(v.addr)
	if err != nil {
		return fmt.Errorf("cannot determine outbound ip: %w", err)
	}

	req.SetHeader("X-Real-IP", ip.String())

	_, err = req.Post(url)
	if err != nil {
		return fmt.Errorf("cannot send HTTP-Request: %w", err)
	}

	return nil
}
