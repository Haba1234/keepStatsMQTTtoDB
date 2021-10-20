package logger

import (
	"testing"

	"github.com/Haba1234/keepStatsMQTTtoDB/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
)

func TestSomething(t *testing.T) {
	l, hook := test.NewNullLogger()
	log, _ := NewLogger(config.LogConf{Level: "debug"})
	log.Logger = l

	log.Error("Hello error")

	require.Equal(t, 1, len(hook.Entries))
	require.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	require.Equal(t, "Hello error", hook.LastEntry().Message)

	hook.Reset()
	require.Nil(t, hook.LastEntry())
}
