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

func (log Logger) WithPrefix(prefix string) *Logger {
	log.prefix = prefix
	return &log
}

// ATTENTION: You SHOULD be very careful about concurrency safety or deadlocks with
// dynamic contexts.
func (log Logger) WithContexts(kvs ...interface{}) *Logger {
	log.staticContexts, log.dynamicContexts =
		appendContexts(log.staticContexts, log.dynamicContexts, kvs)
	return &log
}

func (log Logger) WithMark(ok bool) *Logger {
	log.mark = ok
	return &log
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
