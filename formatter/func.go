// Package formatter provides wrappers to the interface iface.Formatter.
package formatter

import (
	"github.com/gratonos/gxlog/iface"
)

// Do NOT call any method of the Logger within the function, or it may deadlock.
type Func func(record *iface.Record) []byte

func (fn Func) Format(record *iface.Record) []byte {
	return fn(record)
}
