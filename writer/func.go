// Package writer provides wrappers to the interface iface.Writer.
package writer

import (
	"github.com/gratonos/gxlog/iface"
)

// Do NOT call any method of the Logger within the function, or it may deadlock.
type Func func(bs []byte, record *iface.Record) error

func (self Func) Write(bs []byte, record *iface.Record) error {
	return self(bs, record)
}
