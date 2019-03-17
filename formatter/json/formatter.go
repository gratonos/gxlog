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

func New(fileSegs, pkgSegs, funcSegs int) *Formatter {
	formatter := &Formatter{
		fileSegs: fileSegs,
		pkgSegs:  pkgSegs,
		funcSegs: funcSegs,
	}
	return formatter
}

func (formatter *Formatter) Format(record *iface.Record) []byte {
	formatter.lock.Lock()

	buf := formatter.buf[:0]
	buf = append(buf, "{"...)
	buf = formatStrField(buf, "", "Time", record.Time.Format(time.RFC3339Nano), false)
	buf = formatIntField(buf, ",", "Level", int(record.Level))
	file := util.LastSegments(record.File, formatter.fileSegs, '/')
	buf = formatStrField(buf, ",", "File", file, true)
	buf = formatIntField(buf, ",", "Line", record.Line)
	pkg := util.LastSegments(record.Pkg, formatter.pkgSegs, '/')
	buf = formatStrField(buf, ",", "Pkg", pkg, false)
	fn := util.LastSegments(record.Func, formatter.funcSegs, '.')
	buf = formatStrField(buf, ",", "Func", fn, false)
	buf = formatStrField(buf, ",", "Msg", record.Msg, true)
	buf = formatter.formatAux(buf, &record.Aux)
	buf = append(buf, "}\n"...)
	formatter.buf = buf

	formatter.lock.Unlock()

	return buf
}

func (formatter *Formatter) FileSegs() int {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.fileSegs
}

func (formatter *Formatter) SetFileSegs(segs int) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.fileSegs = segs
}

func (formatter *Formatter) PkgSegs() int {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.pkgSegs
}

func (formatter *Formatter) SetPkgSegs(segs int) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.pkgSegs = segs
}

func (formatter *Formatter) FuncSegs() int {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.funcSegs
}

func (formatter *Formatter) SetFuncSegs(segs int) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.funcSegs = segs
}

func (formatter *Formatter) formatAux(buf []byte, aux *iface.Auxiliary) []byte {
	buf = append(buf, `,"Aux":{`...)
	buf = formatStrField(buf, "", "Prefix", aux.Prefix, true)
	buf = formatContexts(buf, aux.Contexts)
	buf = formatBoolField(buf, ",", "Mark", aux.Mark)
	return append(buf, "}"...)
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
	}
	return append(buf, "false"...)
}
