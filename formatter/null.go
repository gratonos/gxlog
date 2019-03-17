package formatter

import (
	"github.com/gratonos/gxlog/iface"
)

var nullFormatter = Func(func(*iface.Record) []byte { return nil })

func Null() iface.Formatter {
	return nullFormatter
}
