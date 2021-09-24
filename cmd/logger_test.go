package cmd

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type bufferWriteSyncer struct {
	*bytes.Buffer
}

func (bws bufferWriteSyncer) Sync() error {
	return nil
}

func newTestErrorLogger(verbose, useStructuredLogging bool) (*zap.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	syncer := bufferWriteSyncer{Buffer: buf}
	logger := newErrorLogger(verbose, useStructuredLogging, syncer)
	return logger, syncer.Buffer
}

func newTestOutputLogger() (*zap.Logger, *bytes.Buffer) {
	buf := new(bytes.Buffer)
	syncer := bufferWriteSyncer{Buffer: buf}
	logger := newOutputLogger(syncer)
	return logger, syncer.Buffer
}

func TestErrorLoggerSchema(t *testing.T) {
	// Expected schema is:
	// {
	//     "version": "vX.X.X",
	//     "timestamp": "RFC 3339-encoded timestamp",
	//     "error_details": {
	//         "error": "error type",
	//         "context": "additional error context (optional)"
	//     }
	// }
	verbose := true
	useStructuredLogging := true

	errorLogger, buf := newTestErrorLogger(verbose, useStructuredLogging)
	errorLogger.Info("This is a message")

	outputJSON := buf.Bytes()

	var output map[string]interface{}
	json.Unmarshal(outputJSON, &output)

	expectedOutput := map[string]interface{}{
		"version":   "v1.0.0",
		"timestamp": "tested separately",
		"error_details": map[string]interface{}{
			"error":   "This is a message",
			"context": "",
		},
	}

	a := assert.New(t)
	a.Equal(len(output), len(expectedOutput))

	for k := range expectedOutput {
		a.Contains(output, k)
	}

	a.Equal(output["version"], expectedOutput["version"])
	a.Equal(output["error_details"], expectedOutput["error_details"])

	// This should not fail
	_, err := time.Parse(time.RFC3339Nano, output["timestamp"].(string))
	a.Nil(err)
}
