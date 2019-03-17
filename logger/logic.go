package logger

import (
	"github.com/gratonos/gxlog/iface"
)

func And(filter Filter, filters ...Filter) Filter {
	return func(record *iface.Record) bool {
		if !filter(record) {
			return false
		}
		for _, filter := range filters {
			if !filter(record) {
				return false
			}
		}
		return true
	}
}

func Or(filter Filter, filters ...Filter) Filter {
	return func(record *iface.Record) bool {
		if filter(record) {
			return true
		}
		for _, filter := range filters {
			if filter(record) {
				return true
			}
		}
		return false
	}
}

func Not(filter Filter) Filter {
	return func(record *iface.Record) bool {
		return !filter(record)
	}
}
