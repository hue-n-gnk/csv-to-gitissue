package logger

import (
	"context"
	"fmt"
	"time"

	gormlogger "gorm.io/gorm/logger"
	gormutils "gorm.io/gorm/utils"
)

// Customize SQL Logger for gorm library
// ref: https://github.com/wantedly/gorm-zap
// ref: https://github.com/go-gorm/gorm/blob/master/logger/logger.go

// Logger is an alternative implementation of *gorm.Logger
type GormLogger interface {
	LogMode(level gormlogger.LogLevel) gormlogger.Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error)
}

type gormLogger struct {
	logger Logger
}

func NewGormLogger(logger Logger) GormLogger {
	return &gormLogger{
		logger: logger,
	}
}

const (
	logTitle      = "[gorm] "
	sqlFormat     = logTitle + "%s"
	messageFormat = logTitle + "%s, %s"
	errorFormat   = logTitle + "%s, %s, %s"
	slowThreshold = 200
)

// LogMode The log level of gorm logger is overwrited by the log level of Zap logger.
func (l *gormLogger) LogMode(_ gormlogger.LogLevel) gormlogger.Interface {
	return l
}

// Info prints a information log.
func (l *gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Infof(messageFormat, append([]interface{}{msg, gormutils.FileWithLineNum()}, data...)...)
}

// Warn prints a warning log.
func (l *gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Warnf(messageFormat, append([]interface{}{msg, gormutils.FileWithLineNum()}, data...)...)
}

// Error prints a error log.
func (l *gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Warnf(messageFormat, append([]interface{}{msg, gormutils.FileWithLineNum()}, data...)...)
}

// Trace prints a trace log such as sql, source file and error.
func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil:
		sql, _ := fc()
		l.logger.WithContext(ctx).Warnf(errorFormat, gormutils.FileWithLineNum(), err, sql)
	case elapsed > slowThreshold*time.Millisecond && slowThreshold*time.Millisecond != 0:
		sql, _ := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", slowThreshold)
		l.logger.WithContext(ctx).Warnf(errorFormat, gormutils.FileWithLineNum(), slowLog, sql)
	default:
		sql, _ := fc()
		l.logger.WithContext(ctx).Infof(sqlFormat, sql)
	}
}
