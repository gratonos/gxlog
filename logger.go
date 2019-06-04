// Package gxlog provides the default Logger and the default Formatter.
package gxlog

import (
	"os"

	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/logger"
	"github.com/gratonos/gxlog/writer"
)

var (
	defaultLogger    *logger.Logger
	defaultFormatter *text.Formatter
)

func init() {
	defaultLogger = logger.New(logger.Config{})
	defaultFormatter = text.New(text.Config{Coloring: true})

	defaultLogger.SetSlot(logger.Slot0, logger.Slot{
		Formatter: defaultFormatter,
		Writer:    writer.Wrap(os.Stderr),
	})
}

func Logger() *logger.Logger {
	return defaultLogger
}

func Formatter() *text.Formatter {
	return defaultFormatter
}
