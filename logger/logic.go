package logger

import (
	"github.com/gratonos/gxlog/iface"
)

func And(filters ...Filter) Filter {
	fillFilters(filters)

	return func(record *iface.Record) bool {
		for _, filter := range filters {
			if !filter(record) {
				return false
			}
		}
		return true
	}
}

func Or(filters ...Filter) Filter {
	fillFilters(filters)

	return func(record *iface.Record) bool {
		for _, filter := range filters {
			if filter(record) {
				return true
			}
		}
		return false
	}
}

func Not(filter Filter) Filter {
	if filter == nil {
		filter = nullFilter
	}
	return func(record *iface.Record) bool {
		return !filter(record)
	}
}

func fillFilters(filters []Filter) {
	for i, filter := range filters {
		if filter == nil {
			filters[i] = nullFilter
		}
	}
}
