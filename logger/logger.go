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

type additional struct {
	Level    iface.Level
	Filter   Filter
	Prefix   string
	Statics  []iface.Context
	Dynamics []dynamicContext
	Mark     bool
}

// A Logger is a logging framework that contains EIGHT slots. Each Slot contains
// a Formatter and a Writer. A Logger has its own level and filter while each
// Slot has its independent level and filter. Logger calls the Formatter and
// Writer of each Slot in the order from Slot0 to Slot7 when a log is emitted.
type Logger struct {
	additional additional // copy on write, concurrency safe

	config      *Config
	slots       []Slot
	equivalents [][]int // indexes of equivalent formatters
	lock        *sync.Mutex
}

func New(config Config) *Logger {
	config.SetDefaults()
	return &Logger{
		additional: additional{
			Filter: nullFilter,
		},
		config:      &config,
		slots:       initSlots(),
		equivalents: make([][]int, MaxSlot),
		lock:        new(sync.Mutex),
	}
}

func (this *Logger) Level() iface.Level {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Level
}

func (this *Logger) SetLevel(level iface.Level) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Level = level
}

func (this *Logger) Filter() Filter {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.config.Filter
}

func (this *Logger) SetFilter(filter Filter) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.config.Filter = fillFilter(filter)
}

func (this *Logger) Trace(args ...interface{}) {
	this.Log(1, iface.Trace, args...)
}

func (this *Logger) Tracef(fmtstr string, args ...interface{}) {
	this.Logf(1, iface.Trace, fmtstr, args...)
}

func (this *Logger) Debug(args ...interface{}) {
	this.Log(1, iface.Debug, args...)
}

func (this *Logger) Debugf(fmtstr string, args ...interface{}) {
	this.Logf(1, iface.Debug, fmtstr, args...)
}

func (this *Logger) Info(args ...interface{}) {
	this.Log(1, iface.Info, args...)
}

func (this *Logger) Infof(fmtstr string, args ...interface{}) {
	this.Logf(1, iface.Info, fmtstr, args...)
}

func (this *Logger) Warn(args ...interface{}) {
	this.Log(1, iface.Warn, args...)
}

func (this *Logger) Warnf(fmtstr string, args ...interface{}) {
	this.Logf(1, iface.Warn, fmtstr, args...)
}

func (this *Logger) Error(args ...interface{}) {
	this.Log(1, iface.Error, args...)
}

func (this *Logger) Errorf(fmtstr string, args ...interface{}) {
	this.Logf(1, iface.Error, fmtstr, args...)
}

func (this *Logger) Fatal(args ...interface{}) {
	this.Log(1, iface.Fatal, args...)
	os.Exit(1)
}

func (this *Logger) Fatalf(fmtstr string, args ...interface{}) {
	this.Logf(1, iface.Fatal, fmtstr, args...)
	os.Exit(1)
}

func (this *Logger) EError(args ...interface{}) error {
	msg := fmt.Sprint(args...)
	this.Log(1, iface.Error, msg)
	return errors.New(msg)
}

func (this *Logger) EErrorf(fmtstr string, args ...interface{}) error {
	msg := fmt.Sprintf(fmtstr, args...)
	this.Log(1, iface.Error, msg)
	return errors.New(msg)
}

func (this *Logger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	this.Log(1, iface.Fatal, msg)
	panic(msg)
}

func (this *Logger) Panicf(fmtstr string, args ...interface{}) {
	msg := fmt.Sprintf(fmtstr, args...)
	this.Log(1, iface.Fatal, msg)
	panic(msg)
}

func (this *Logger) Log(callDepth int, level iface.Level, args ...interface{}) {
	if this.needToLog(level) {
		this.log(callDepth, level, fmt.Sprint(args...))
	}
}

func (this *Logger) Logf(callDepth int, level iface.Level, fmtstr string, args ...interface{}) {
	if this.needToLog(level) {
		this.log(callDepth, level, fmt.Sprintf(fmtstr, args...))
	}
}

func (this *Logger) Timing(level iface.Level, args ...interface{}) func() {
	if this.needToLog(level) {
		return this.doneFunc(level, fmt.Sprint(args...))
	} else {
		return func() {}
	}
}

func (this *Logger) Timingf(level iface.Level, fmtstr string, args ...interface{}) func() {
	if this.needToLog(level) {
		return this.doneFunc(level, fmt.Sprintf(fmtstr, args...))
	} else {
		return func() {}
	}
}

func (this *Logger) needToLog(level iface.Level) (need bool) {
	if level < iface.Trace || level > iface.Fatal {
		panic(fmt.Sprintf("gxlog: invalid log level: %d", level))
	}

	this.lock.Lock()
	if this.additional.Level <= level && this.config.Level <= level {
		need = true
	}
	this.lock.Unlock()

	return need
}

func (this *Logger) doneFunc(level iface.Level, msg string) func() {
	now := time.Now()
	return func() {
		cost := time.Since(now)
		this.log(0, level, fmt.Sprintf("%s (cost: %v)", msg, cost))
	}
}

func (this *Logger) log(callDepth int, level iface.Level, msg string) {
	file, line, pkg, fn := getPosInfo(callDepth + callDepthOffset)

	contexts := this.additional.Statics
	for _, context := range this.additional.Dynamics {
		contexts = append(contexts, iface.Context{
			Key:   fmt.Sprint(context.Key),
			Value: fmt.Sprint(context.Value(context.Key)),
		})
	}

	this.lock.Lock()

	record := &iface.Record{
		Time:     time.Now(),
		Level:    level,
		File:     file,
		Line:     line,
		Pkg:      pkg,
		Func:     fn,
		Msg:      msg,
		Prefix:   this.additional.Prefix,
		Contexts: contexts,
		Mark:     this.additional.Mark,
	}

	if this.additional.Filter(record) && this.config.Filter(record) {
		this.formatAndWrite(level, record)
	}

	this.lock.Unlock()
}

func (this *Logger) formatAndWrite(level iface.Level, record *iface.Record) {
	var formats [MaxSlot][]byte
	for i := 0; i < MaxSlot; i++ {
		slot := &this.slots[i]
		if slot.Level > level || !slot.Filter(record) {
			continue
		}

		format := formats[i]
		if format == nil {
			format = slot.Formatter.Format(record)
			for _, id := range this.equivalents[i] {
				formats[id] = format
			}
		}

		err := slot.Writer.Write(format, record)
		if err != nil {
			slot.ErrorHandler(format, record, err)
		}
	}
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
