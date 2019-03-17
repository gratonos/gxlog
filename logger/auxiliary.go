package logger

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

type Dynamic func(key interface{}) interface{}

type dynamicContext struct {
	Key   interface{}
	Value Dynamic
}

func (log *Logger) Prefix(prefix string) *Logger {
	clone := *log
	clone.prefix = prefix
	return &clone
}

// ATTENTION: You SHOULD be very careful about concurrency safety or deadlocks with
// dynamic contexts.
func (log *Logger) Contexts(kvs ...interface{}) *Logger {
	clone := *log
	clone.staticContexts, clone.dynamicContexts =
		appendContexts(clone.staticContexts, clone.dynamicContexts, kvs)
	return &clone
}

func (log *Logger) Mark(ok bool) *Logger {
	clone := *log
	clone.mark = ok
	return &clone
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
	// slicing to set the capacity of slice to its length, force the next
	//   appending to the slice to reallocate memory
	return staticContexts[:len(staticContexts):len(staticContexts)],
		dynamicContexts[:len(dynamicContexts):len(dynamicContexts)]
}
