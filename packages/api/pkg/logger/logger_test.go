package logger

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Parallel()

	log, err := NewLogger()
	require.NoError(t, err)
	{
		log.Info("test log info")
	}
	{
		log.Infof("test log infof")
	}
	{
		log.Debug("test log debug")
	}
	{
		log.Debugf("test log debugf")
	}
	{
		log.Warn("test log warn")
	}
	{
		log.Warnf("test log warnf")
	}
	{
		log.Error("test log error")
	}
	{
		log.Errorf("test log errorf")
	}
}

func TestSync(t *testing.T) {
	t.Parallel()

	{
		logg, err := NewLogger()
		require.NoError(t, err)
		err = logg.Sync()
		require.Error(t, err)
	}

	{
		temp := struct{ *logger }{logger: nil}
		err := temp.Sync()
		require.NoError(t, err)
	}
}

func TestLoggerWithCustomMessage(t *testing.T) {
	t.Parallel()

	log, err := NewLogger()
	require.NoError(t, err)
	{
		log.Info("test log info", FieldMap{
			"message": nil,
		})
	}
	{
		log.Info("test log info", FieldMap{
			"message": "test",
		})
	}
	{
		log.Error("test error", FieldMap{
			"error": fmt.Errorf("failed"),
		})
	}
	{
		log.Error("test error", FieldMap{
			"error": nil,
		})
	}
	{
		log.Error("", FieldMap{
			"error": fmt.Errorf("failed"),
		})
	}
	{
		log.Error("test error", FieldMap{
			"error": errors.New(""),
		})
	}
}

func TestErrorToString(t *testing.T) {
	t.Parallel()

	{
		res := errorToString(nil, 2)
		require.NotNil(t, res)
		require.Equal(t, "<nil>", res)
	}

	{
		err := fmt.Errorf(`
			failed
			runtime failed
			failed1
			failed2
		`)
		res := errorToString(err, 2)
		require.NotNil(t, res)
	}
	{
		err := fmt.Errorf("failed")
		res := errorToString(err, 0)
		require.NotNil(t, res)
		require.Equal(t, res, "failed")
	}
}
