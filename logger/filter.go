package logger

import (
	"github.com/gratonos/gxlog/iface"
)

// Do NOT call any method of the Logger within a Filter, or it may deadlock.
type Filter func(record *iface.Record) bool

var nullFilter = func(*iface.Record) bool { return true }

func NullFilter() Filter {
	return nullFilter
}
