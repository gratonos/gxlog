package logger

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

// ATTENTION: You SHOULD be careful about concurrency safety or deadlocks with dynamic contexts.
type Dynamic func(key interface{}) (value interface{})

type dynamicContext struct {
	Key   interface{}
	Value Dynamic
}

func (self Logger) WithLevel(level iface.Level) *Logger {
	self.additional.Level = level
	return &self
}

func (self Logger) WithFilter(filter Filter) *Logger {
	self.additional.Filter = fillFilter(filter)
	return &self
}

func (self Logger) WithPrefix(prefix string) *Logger {
	self.additional.Prefix = prefix
	return &self
}

func (self Logger) WithContexts(kvs ...interface{}) *Logger {
	self.additional.Statics, self.additional.Dynamics =
		appendContexts(self.additional.Statics, self.additional.Dynamics, kvs)
	return &self
}

func (self Logger) WithMark(mark bool) *Logger {
	self.additional.Mark = mark
	return &self
}

func appendContexts(staticContexts []iface.Context, dynamicContexts []dynamicContext,
	kvs []interface{}) ([]iface.Context, []dynamicContext) {

	for len(kvs) >= 2 {
		dynamic, ok := kvs[1].(Dynamic)
		if ok {
			dynamicContexts = append(dynamicContexts, dynamicContext{
				Key:   kvs[0],
				Value: dynamic,
			})
		} else {
			staticContexts = append(staticContexts, iface.Context{
				Key:   fmt.Sprint(kvs[0]),
				Value: fmt.Sprint(kvs[1]),
			})
		}
		kvs = kvs[2:]
	}
	// Slicing to set the capacity of slice to its length, force the next appending to the slice
	// to reallocate memory. This ensures all the slices refer to different pieces of memory and
	// avoids data contention.
	return staticContexts[:len(staticContexts):len(staticContexts)],
		dynamicContexts[:len(dynamicContexts):len(dynamicContexts)]
}
