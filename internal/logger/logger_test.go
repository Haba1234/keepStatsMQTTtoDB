package logger

import (
	"testing"

	"github.com/Haba1234/keepStatsMQTTtoDB/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Parallel()
	l, hook := test.NewNullLogger()
	log, _ := NewLogger(config.LogConf{Level: "debug"})
	l.Level = logrus.DebugLevel
	log.Logger = l

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		log.Error("hello error")

		require.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		require.Equal(t, "hello error", hook.LastEntry().Message)
	})

	t.Run("Errorf", func(t *testing.T) {
		t.Parallel()
		log.Errorf("error: %s", "hello")

		require.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
		require.Equal(t, "error: hello", hook.LastEntry().Message)
	})

	t.Run("Info", func(t *testing.T) {
		t.Parallel()
		log.Info("hello info")

		require.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
		require.Equal(t, "hello info", hook.LastEntry().Message)
	})

	t.Run("Infof", func(t *testing.T) {
		t.Parallel()
		log.Infof("info: %s", "hello")

		require.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
		require.Equal(t, "info: hello", hook.LastEntry().Message)
	})

	t.Run("Debug", func(t *testing.T) {
		t.Parallel()
		log.Debug(123)

		require.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
		require.Equal(t, "123", hook.LastEntry().Message)
	})

	t.Run("Debug with parameters", func(t *testing.T) {
		t.Parallel()
		log.Debug("param", 123, "hello", "test")
		require.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
		require.Equal(t, "[hello test]", hook.LastEntry().Message)
	})

	hook.Reset()
	require.Nil(t, hook.LastEntry())
}
