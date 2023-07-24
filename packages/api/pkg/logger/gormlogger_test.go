package logger

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGormLogger(t *testing.T) {

	logger, err := NewLogger()
	require.NoError(t, err)

	logg := NewGormLogger(logger)
	ctx := context.TODO()
	{
		logg.Info(ctx, "test gorm")
	}
	{
		logg.Error(ctx, "test gorm")
	}
	{
		logg.Warn(ctx, "test gorm")
	}
	{
		l := logg.LogMode(1)
		require.NotNil(t, l)
	}
}

func TestTrace(t *testing.T) {
	t.Parallel()
	testFunc := func() (string, int64) {
		return "", 0
	}

	logger, err := NewLogger()
	require.NoError(t, err)

	logg := NewGormLogger(logger)

	ctx := context.TODO()

	{
		begin := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		logg.Trace(ctx, begin, testFunc, fmt.Errorf("failed"))
	}

	{
		begin := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		logg.Trace(ctx, begin, testFunc, nil)
	}

	{
		begin := time.Now().Add(time.Hour)
		logg.Trace(ctx, begin, testFunc, nil)
	}
}
