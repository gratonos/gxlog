// Package logger implements a concise, powerful, flexible and extensible logger.
package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gratonos/gxlog/iface"
)

const (
	callDepthOffset = 3
)

// A Logger is a logging framework that contains EIGHT slots. Each Slot contains
// a Formatter and a Writer. A Logger has its own level and filter while each
// Slot has its independent level and filter. Logger calls the Formatter and
// Writer of each Slot in the order from Slot0 to Slot7 when a log is emitted.
type Logger struct {
	config *Config

	level           iface.Level
	filter          Filter
	prefix          string
	staticContexts  []iface.Context
	dynamicContexts []dynamicContext
	mark            bool

	slots       []Slot
	equivalents [][]int // indexes of equivalent formatters
	lock        *sync.Mutex
}

func New(config Config) *Logger {
	return &Logger{
		config:      &config,
		slots:       initSlots(),
		equivalents: make([][]int, MaxSlot),
		lock:        new(sync.Mutex),
	}
}

func (log *Logger) Level() iface.Level {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Level
}

func (log *Logger) SetLevel(level iface.Level) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.Level = level
}

func (log *Logger) Filter() Filter {
	log.lock.Lock()
	defer log.lock.Unlock()

	return log.config.Filter
}

func (log *Logger) SetFilter(filter Filter) {
	log.lock.Lock()
	defer log.lock.Unlock()

	log.config.Filter = filter
}

func (log *Logger) Trace(args ...interface{}) {
	log.Log(1, iface.Trace, args...)
}

func (log *Logger) Tracef(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Trace, fmtstr, args...)
}

func (log *Logger) Debug(args ...interface{}) {
	log.Log(1, iface.Debug, args...)
}

func (log *Logger) Debugf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Debug, fmtstr, args...)
}

func (log *Logger) Info(args ...interface{}) {
	log.Log(1, iface.Info, args...)
}

func (log *Logger) Infof(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Info, fmtstr, args...)
}

func (log *Logger) Warn(args ...interface{}) {
	log.Log(1, iface.Warn, args...)
}

func (log *Logger) Warnf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Warn, fmtstr, args...)
}

func (log *Logger) Error(args ...interface{}) {
	log.Log(1, iface.Error, args...)
}

func (log *Logger) Errorf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Error, fmtstr, args...)
}

func (log *Logger) Fatal(args ...interface{}) {
	log.Log(1, iface.Fatal, args...)
	os.Exit(1)
}

func (log *Logger) Fatalf(fmtstr string, args ...interface{}) {
	log.Logf(1, iface.Fatal, fmtstr, args...)
	os.Exit(1)
}

func (log *Logger) EError(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	log.Log(1, iface.Error, msg)
	return errors.New(msg)
}

func (log *Logger) EErrorf(fmtstr string, args ...interface{}) error {
	msg := fmt.Sprintf(fmtstr, args...)
	log.Log(1, iface.Error, msg)
	return errors.New(msg)
}

func (log *Logger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	log.Log(1, iface.Fatal, msg)
	panic(msg)
}

func (log *Logger) Panicf(fmtstr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	log.Log(1, iface.Fatal, msg)
	panic(msg)
}

func (log *Logger) Log(callDepth int, level iface.Level, args ...interface{}) {
	checkLevel(level)

	if log.level > level {
		return
	}

	log.lock.Lock()
	logLevel := log.config.Level
	log.lock.Unlock()

	if logLevel <= level {
		log.write(callDepth, level, fmt.Sprint(args...))
	}
}

func (log *Logger) Logf(callDepth int, level iface.Level, fmtstr string, args ...interface{}) {
	checkLevel(level)

	if log.level > level {
		return
	}

	log.lock.Lock()
	logLevel := log.config.Level
	log.lock.Unlock()

	if logLevel <= level {
		log.write(callDepth, level, fmt.Sprintf(fmtstr, args...))
	}
}

func (log *Logger) Timing(level iface.Level, args ...interface{}) func() {
	checkLevel(level)

	if log.level > level {
		return func() {}
	}

	log.lock.Lock()
	logLevel := log.config.Level
	log.lock.Unlock()

	if logLevel <= level {
		return log.doneFunc(level, fmt.Sprint(args...))
	}
	return func() {}
}

func (log *Logger) Timingf(level iface.Level, fmtstr string, args ...interface{}) func() {
	checkLevel(level)

	if log.level > level {
		return func() {}
	}

	log.lock.Lock()
	logLevel := log.config.Level
	log.lock.Unlock()

	if logLevel <= level {
		return log.doneFunc(level, fmt.Sprintf(fmtstr, args...))
	}
	return func() {}
}

func (log *Logger) doneFunc(level iface.Level, msg string) func() {
	now := time.Now()
	return func() {
		cost := time.Since(now)
		log.write(0, level, fmt.Sprintf("%s (cost: %v)", msg, cost))
	}
}

func (log *Logger) write(callDepth int, level iface.Level, msg string) {
	file, line, pkg, fn := getPosInfo(callDepth + callDepthOffset)

	log.lock.Lock()

	record := &iface.Record{
		Time:  time.Now(),
		Level: level,
		File:  file,
		Line:  line,
		Pkg:   pkg,
		Func:  fn,
		Msg:   msg,
	}

	log.attachAuxiliary(record)

	if (log.config.Filter != nil && !log.config.Filter(record)) ||
		(log.filter != nil && !log.filter(record)) {
		log.lock.Unlock()
		return
	}

	var formats [MaxSlot][]byte
	for i := 0; i < MaxSlot; i++ {
		slot := &log.slots[i]
		if slot.Level > level {
			continue
		}
		if slot.Filter != nil && !slot.Filter(record) {
			continue
		}
		format := formats[i]
		if format == nil && slot.Formatter != nil {
			format = slot.Formatter.Format(record)
			for _, id := range log.equivalents[i] {
				formats[id] = format
			}
		}
		if slot.Writer != nil {
			err := slot.Writer.Write(format, record)
			if err != nil && slot.ErrorHandler != nil {
				slot.ErrorHandler(format, record, err)
			}
		}
	}

	log.lock.Unlock()
}

func (log *Logger) attachAuxiliary(record *iface.Record) {
	record.Aux.Prefix = log.prefix
	record.Aux.Contexts = log.staticContexts
	for _, context := range log.dynamicContexts {
		record.Aux.Contexts = append(record.Aux.Contexts, iface.Context{
			Key:   fmt.Sprint(context.Key),
			Value: fmt.Sprint(context.Value(context.Key)),
		})
	}
	record.Aux.Mark = log.mark
}

func getPosInfo(callDepth int) (file string, line int, pkg, fn string) {
	var pc uintptr
	var ok bool
	pc, file, line, ok = runtime.Caller(callDepth)
	if ok {
		name := runtime.FuncForPC(pc).Name()
		pkg, fn = splitPkgAndFunc(name)
	} else {
		file = "?file?"
		line = -1
		pkg = "?pkg?"
		fn = "?func?"
	}
	return filepath.ToSlash(file), line, pkg, fn
}

func splitPkgAndFunc(name string) (string, string) {
	lastSlash := strings.LastIndexByte(name, '/')
	nextDot := strings.IndexByte(name[lastSlash+1:], '.')
	if nextDot < 0 {
		return "?pkg?", "?func?"
	}
	nextDot += (lastSlash + 1)
	return name[:nextDot], name[nextDot+1:]
}

func checkLevel(level iface.Level) {
	if level < iface.Trace || level > iface.Fatal {
		panic(fmt.Sprintf("gxlog: invalid log level: %d", level))
	}
}
