// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)



type errorDetails struct {
	// error message
	error string
	// stacktrace or additional context
	context string
}

func (ed errorDetails) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("error", ed.error)
	enc.AddString("context", ed.context)
	return nil
}

type structuredEncoder struct {
	zapcore.Encoder
}

func (enc structuredEncoder) Clone() zapcore.Encoder {
	return structuredEncoder{
		Encoder: enc.Encoder.Clone(),
	}
}

func (enc structuredEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	ed := errorDetails{
		error: ent.Message,
		context: ent.Stack,
	}
	fields = append([]zapcore.Field{
		zap.Object("error_details", ed),
	}, fields...)
	return enc.Encoder.EncodeEntry(ent, fields)
}

func newJSONEncoder(cfg zapcore.EncoderConfig) structuredEncoder {
	return structuredEncoder{
		Encoder: zapcore.NewJSONEncoder(cfg),
	}
}

func newConsoleEncoder(cfg zapcore.EncoderConfig) structuredEncoder {
	return structuredEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}

func newLogger(verbose, useStructuredLogging bool) *zap.Logger {
	// Return a logger that produces expected structured output format
	var level zap.AtomicLevel
	options := []zap.Option{
		zap.Fields(
			// Message format version
			zap.String("version", "v1.0.0"),
		),
	}

	if verbose {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
		options = append(options, zap.AddStacktrace(zap.WarnLevel))
	} else {
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
		options = append(options, zap.AddStacktrace(zap.ErrorLevel))
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		MessageKey:     "",
		StacktraceKey:  "",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
	}
	var encoder structuredEncoder
	if useStructuredLogging {
		encoder = newJSONEncoder(encoderConfig)
	} else {
		encoder = newConsoleEncoder(encoderConfig)
	}
	core := zapcore.NewCore(encoder, zapcore.Lock(os.Stderr), level)
	return zap.New(core, options...)
}
