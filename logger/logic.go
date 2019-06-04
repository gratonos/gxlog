package logger

import (
	"github.com/gratonos/gxlog/iface"
)

func And(filters ...Filter) Filter {
	return func(record *iface.Record) bool {
		for _, filter := range filters {
			if filter != nil && !filter(record) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...Filter) Filter {
	return func(record *iface.Record) bool {
		for _, filter := range filters {
			if filter == nil || filter(record) {
				return true
			}
		}
		return false
	}
}

func Not(filter Filter) Filter {
	if filter == nil {
		return func(*iface.Record) bool {
			return false
		}
	}
	return func(record *iface.Record) bool {
		return !filter(record)
	}
}
