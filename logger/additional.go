package logger

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

// ATTENTION: You SHOULD be careful about concurrency safety or deadlocks with dynamic contexts.
type Dynamic func() interface{}

type dynamicContext struct {
	Key   string
	Value Dynamic
}

type additional struct {
	Level    iface.Level
	Filter   Filter
	Prefix   string
	Statics  []iface.Context
	Dynamics []dynamicContext
	Mark     bool
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

func (self Logger) WithStatics(kvs ...interface{}) *Logger {
	statics := cloneStatics(self.additional.Statics, keyValuePairs(len(kvs)))
	self.additional.Statics = appendStatics(statics, kvs)
	return &self
}

func (self Logger) WithDynamics(kvs ...interface{}) *Logger {
	dynamics := cloneDynamics(self.additional.Dynamics, keyValuePairs(len(kvs)))
	self.additional.Dynamics = appendDynamics(dynamics, kvs)
	return &self
}

func (self Logger) WithMark(mark bool) *Logger {
	self.additional.Mark = mark
	return &self
}

func keyValuePairs(n int) int {
	if n%2 != 0 {
		panic("gxlog: unmatched key/value pair")
	}
	return n / 2
}

func cloneStatics(statics []iface.Context, delta int) []iface.Context {
	clone := make([]iface.Context, len(statics), len(statics)+delta)
	copy(clone, statics)
	return clone
}

func appendStatics(statics []iface.Context, kvs []interface{}) []iface.Context {
	// len(kvs) is even, checked by func keyValuePairs
	for len(kvs) > 0 {
		statics = append(statics, iface.Context{
			Key:   fmt.Sprint(kvs[0]),
			Value: fmt.Sprint(kvs[1]),
		})
		kvs = kvs[2:]
	}
	return statics
}

func cloneDynamics(dynamics []dynamicContext, delta int) []dynamicContext {
	clone := make([]dynamicContext, len(dynamics), len(dynamics)+delta)
	copy(clone, dynamics)
	return clone
}

func appendDynamics(dynamics []dynamicContext, kvs []interface{}) []dynamicContext {
	// len(kvs) is even, checked by func keyValuePairs
	for len(kvs) > 0 {
		dynamic, ok := kvs[1].(Dynamic)
		if !ok {
			panic("gxlog: dynamic value must be type Dynamic")
		} else if dynamic == nil {
			panic("gxlog: dynamic value must not be nil")
		}
		dynamics = append(dynamics, dynamicContext{
			Key:   fmt.Sprint(kvs[0]),
			Value: dynamic,
		})
		kvs = kvs[2:]
	}
	return dynamics
}
