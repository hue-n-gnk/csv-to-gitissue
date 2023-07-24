package logger

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	RequestID = `request_id`

	// TimeKey represents structured log's "time" key.
	TimeKey = `time`

	// MessageKey represents structured log's "message" key
	MessageKey = `message`

	// SeverityKey represents structured log's "severity" key.
	SeverityKey = `severity`

	// CallerKey represents structured log's "caller" key.
	CallerKey = `caller`

	// ErrorKey represents structured log's "error" key.
	ErrorKey = `error`

	// EnvKey represents structured log's "env" key.
	EnvKey = `env`

	// MaxErrLines defines the maximum number of stack traces to output to the log
	MaxErrLines = 64
)

var commonOptions = []zap.Option{
	zap.AddCallerSkip(1),
}

var apiLoggerIntance Logger
var once sync.Once

// FieldMap represents log entry.
type FieldMap map[string]interface{}

func (f FieldMap) build() (message string, fields []zap.Field) {
	var errorString string
	for k, v := range f {
		if err, ok := v.(error); ok && k == ErrorKey {
			errorString = errorToString(err, MaxErrLines)
		} else if k == MessageKey {
			if v == nil {
				message = ""
			} else {
				message = v.(string)
			}
		} else {
			fields = append(fields, zap.Any(k, v))
		}
	}
	if errorString != "" {
		if message == "" {
			message = "error "
		}
		message += " -- " + errorString
	}
	return
}

type Logger interface {
	Debug(message string, f ...FieldMap)
	Debugf(format string, a ...interface{})
	Info(message string, f ...FieldMap)
	Infof(format string, a ...interface{})
	Warn(message string, f ...FieldMap)
	Warnf(format string, a ...interface{})
	Error(message string, f ...FieldMap)
	Errorf(format string, a ...interface{})
	Sync() error
	WithContext(ctx context.Context) Logger
}

type logger struct {
	logger *zap.Logger
}

func NewLogger(params ...interface{}) (Logger, error) {
	level := zap.NewAtomicLevel()
	dev := false
	if os.Getenv("ENV") == "develop" || os.Getenv("ENV") == "development" {
		level.SetLevel(zap.DebugLevel)
		dev = true
	}
	zapConfig := zap.Config{
		Level:       level,
		Development: dev,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			LineEnding:     zapcore.DefaultLineEnding,
			LevelKey:       SeverityKey,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			TimeKey:        TimeKey,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			CallerKey:      CallerKey,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			MessageKey:     MessageKey,
		},
		DisableCaller:    false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if len(params) > 0 {
		zapConfig.InitialFields = params[0].(map[string]interface{})
		zapConfig.InitialFields[EnvKey] = os.Getenv("ENV")
	}

	l, err := zapConfig.Build()
	if err != nil {
		return nil, errors.Wrap(err, "error building zap config")
	}

	logger := &logger{
		logger: l.WithOptions(commonOptions...),
	}

	return logger, nil
}

func MergeToFieldMap(message string, f []FieldMap) *FieldMap {
	fieldMap := FieldMap{}
	fieldMap[MessageKey] = message
	for _, element := range f {
		for key, value := range element {
			fieldMap[key] = value
		}
	}
	return &fieldMap
}

// Debug outputs message with DEBUG severity.
func (l *logger) Debug(message string, f ...FieldMap) {
	if l != nil {
		message, fields := MergeToFieldMap(message, f).build()
		l.logger.Debug(message, fields...)
	}
}

// Debugf outputs formatted message with DEBUG severity.
func (l *logger) Debugf(format string, a ...interface{}) {
	if l != nil {
		l.logger.Debug(fmt.Sprintf(format, a...))
	}
}

// Info outputs message with INFO severity.
func (l *logger) Info(message string, f ...FieldMap) {
	if l != nil {
		message, fields := MergeToFieldMap(message, f).build()
		l.logger.Info(message, fields...)
	}
}

// Infof outputs formatted message with INFO severity.
func (l *logger) Infof(format string, a ...interface{}) {
	if l != nil {
		l.logger.Info(fmt.Sprintf(format, a...))
	}
}

// Warn outputs message with WARNING severity.
func (l *logger) Warn(message string, f ...FieldMap) {
	if l != nil {
		message, fields := MergeToFieldMap(message, f).build()
		l.logger.Warn(message, fields...)
	}
}

// Warnf outputs formatted message with WARNING severity.
func (l *logger) Warnf(format string, a ...interface{}) {
	if l != nil {
		l.logger.Warn(fmt.Sprintf(format, a...))
	}
}

// Error outputs message with ERROR severity.
func (l *logger) Error(message string, f ...FieldMap) {
	if l != nil {
		message, fields := MergeToFieldMap(message, f).build()
		l.logger.Error(message, fields...)
	}
}

// Errorf outputs formatted message with ERROR severity.
func (l *logger) Errorf(format string, a ...interface{}) {
	if l != nil {
		l.logger.Error(fmt.Sprintf(format, a...))
	}
}

// Sync call flushes the buffer. Application should call Sync before exiting
func (l *logger) Sync() error {
	if l != nil {
		return l.logger.Sync()
	}
	return nil
}

func errorToString(err error, maxLines int) string {
	if err == nil {
		return "<nil>"
	}

	lines := strings.Split(fmt.Sprintf("%+v", err), "\n")

	if len(lines) < 2 {
		return lines[0]
	}

	lines = FilterStackTrace(lines)
	max := len(lines)

	if max > maxLines {
		max = maxLines
	}

	return strings.Join(lines[:max], " ")
}

func FilterStackTrace(lines []string) []string {
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		// skip runtime error
		if strings.Contains(line, "runtime") {
			continue
		}
		filtered = append(filtered, line)
	}
	return filtered
}

func (l *logger) WithContext(ctx context.Context) Logger {
	var childLogger *zap.Logger
	if middleware.GetReqID(ctx) != "" {
		childLogger = l.logger.With(zap.String(RequestID, middleware.GetReqID(ctx)))
	} else {
		childLogger = l.logger
	}
	return &logger{
		logger: childLogger,
	}
}

func NewAPILogger() Logger {
	once.Do(func() {
		var err error
		apiLoggerIntance, err = NewLogger()
		if err != nil {
			panic(err)
		}
	})
	return apiLoggerIntance
}
