// Package json implements a json formatter.
package json

import (
	"strconv"
	"sync"
	"time"

	"github.com/gratonos/gxlog/formatter/internal/util"
	"github.com/gratonos/gxlog/iface"
)

type Formatter struct {
	fileSegs int
	pkgSegs  int
	funcSegs int

	buf  []byte
	lock sync.Mutex
}

func New(config Config) *Formatter {
	formatter := &Formatter{
		fileSegs: config.FileSegs,
		pkgSegs:  config.PkgSegs,
		funcSegs: config.FuncSegs,
	}
	return formatter
}

func (this *Formatter) Format(record *iface.Record) []byte {
	this.lock.Lock()

	buf := this.buf[:0]
	buf = append(buf, "{"...)

	buf = formatStrField(buf, "", "Time", record.Time.Format(time.RFC3339Nano), false)
	buf = formatIntField(buf, ",", "Level", int(record.Level))
	file := util.LastSegments(record.File, this.fileSegs, '/')
	buf = formatStrField(buf, ",", "File", file, true)
	buf = formatIntField(buf, ",", "Line", record.Line)
	pkg := util.LastSegments(record.Pkg, this.pkgSegs, '/')
	buf = formatStrField(buf, ",", "Pkg", pkg, false)
	fn := util.LastSegments(record.Func, this.funcSegs, '.')
	buf = formatStrField(buf, ",", "Func", fn, false)
	buf = formatStrField(buf, ",", "Msg", record.Msg, true)

	buf = formatStrField(buf, ",", "Prefix", record.Prefix, true)
	buf = formatContexts(buf, record.Contexts)
	buf = formatBoolField(buf, ",", "Mark", record.Mark)

	buf = append(buf, "}\n"...)
	this.buf = buf

	this.lock.Unlock()

	return buf
}

func (this *Formatter) FileSegs() int {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.fileSegs
}

func (this *Formatter) SetFileSegs(segs int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.fileSegs = segs
}

func (this *Formatter) PkgSegs() int {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.pkgSegs
}

func (this *Formatter) SetPkgSegs(segs int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.pkgSegs = segs
}

func (this *Formatter) FuncSegs() int {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.funcSegs
}

func (this *Formatter) SetFuncSegs(segs int) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.funcSegs = segs
}

func formatContexts(buf []byte, contexts []iface.Context) []byte {
	buf = append(buf, `,"Contexts":[`...)
	sep := ""
	for _, context := range contexts {
		buf = append(buf, sep...)
		buf = append(buf, "{"...)
		buf = formatStrField(buf, "", "Key", context.Key, true)
		buf = formatStrField(buf, ",", "Value", context.Value, true)
		buf = append(buf, "}"...)
		sep = ","
	}
	return append(buf, "]"...)
}

func formatStrField(buf []byte, sep, key, value string, esc bool) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":"`...)
	if esc {
		buf = escape(buf, value)
	} else {
		buf = append(buf, value...)
	}
	return append(buf, `"`...)
}

func formatIntField(buf []byte, sep, key string, value int) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":`...)
	return strconv.AppendInt(buf, int64(value), 10)
}

func formatBoolField(buf []byte, sep, key string, value bool) []byte {
	buf = append(buf, sep...)
	buf = append(buf, `"`...)
	buf = append(buf, key...)
	buf = append(buf, `":`...)
	if value {
		return append(buf, "true"...)
	} else {
		return append(buf, "false"...)
	}
}
