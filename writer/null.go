package writer

import (
	"github.com/gratonos/gxlog/iface"
)

var nullWriter = Func(func([]byte, *iface.Record) error { return nil })

func Null() iface.Writer {
	return nullWriter
}
