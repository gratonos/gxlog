// Package gxlog provides the default Logger and the default Formatter.
package gxlog

import (
	"os"

	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/logger"
	"github.com/gratonos/gxlog/writer"
)

var (
	defaultLogger    *logger.Logger
	defaultFormatter *text.Formatter
)

func init() {
	defaultLogger = logger.New(iface.Trace, nil)
	defaultFormatter = text.New(text.CompactHeader, true, nil)
	defaultLogger.Link(logger.Slot0, defaultFormatter, writer.Wrap(os.Stderr),
		iface.Trace, nil, nil)
}

func Logger() *logger.Logger {
	return defaultLogger
}

func Formatter() *text.Formatter {
	return defaultFormatter
}
