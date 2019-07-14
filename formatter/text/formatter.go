// Package text implements a text formatter.
package text

import (
	"regexp"
	"strings"
	"sync"

	"github.com/gratonos/gxlog/iface"
)

var headerRegexp = regexp.MustCompile("{{([^:%]*?)(?::([^%]*?))?(%.*?)?}}")

type Formatter struct {
	header   string
	coloring bool

	colorMgr  *colorMgr
	appenders []*headerAppender
	suffix    string
	buf       []byte
	lock      sync.Mutex
}

func New(config Config) *Formatter {
	config.SetDefaults()
	formatter := &Formatter{
		coloring: config.Coloring,
		colorMgr: newColorMgr(config.ColorMap, config.MarkColor),
	}
	formatter.SetHeader(config.Header)
	return formatter
}

func (this *Formatter) Header() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.header
}

func (this *Formatter) SetHeader(header string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.header = header
	this.appenders = this.appenders[:0]
	var staticText string
	for header != "" {
		indexes := headerRegexp.FindStringSubmatchIndex(header)
		if indexes == nil {
			break
		}
		begin, end := indexes[0], indexes[1]
		staticText += header[:begin]
		element, property, fmtspec := extractElement(header, indexes)
		if this.addAppender(element, property, fmtspec, staticText) {
			staticText = ""
		}
		header = header[end:]
	}
	this.suffix = staticText + header
}

func (this *Formatter) Coloring() bool {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.coloring
}

func (this *Formatter) SetColoring(ok bool) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.coloring = ok
}

func (this *Formatter) Color(level iface.Level) Color {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.colorMgr.Color(level)
}

func (this *Formatter) SetColor(level iface.Level, color Color) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.colorMgr.SetColor(level, color)
}

func (this *Formatter) MapColors(colorMap map[iface.Level]Color) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.colorMgr.MapColors(colorMap)
}

func (this *Formatter) MarkColor() Color {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.colorMgr.MarkColor()
}

func (this *Formatter) SetMarkColor(color Color) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.colorMgr.SetMarkColor(color)
}

func (this *Formatter) Format(record *iface.Record) []byte {
	this.lock.Lock()

	var left, right []byte
	if this.coloring {
		if record.Mark {
			left, right = this.colorMgr.MarkColorEars()
		} else {
			left, right = this.colorMgr.ColorEars(record.Level)
		}
	}

	buf := this.buf[:0]
	buf = append(buf, left...)
	for _, appender := range this.appenders {
		buf = appender.AppendHeader(buf, record)
	}
	buf = append(buf, this.suffix...)
	buf = append(buf, right...)
	this.buf = buf

	this.lock.Unlock()
	return buf
}

func (this *Formatter) addAppender(element, property, fmtspec, staticText string) bool {
	appender := newHeaderAppender(element, property, fmtspec, staticText)
	if appender == nil {
		return false
	}
	this.appenders = append(this.appenders, appender)
	return true
}

func extractElement(header string, indexes []int) (element, property, fmtspec string) {
	element = strings.ToLower(getField(header, indexes[2], indexes[3]))
	property = getField(header, indexes[4], indexes[5])
	fmtspec = getField(header, indexes[6], indexes[7])
	if fmtspec == "%" {
		fmtspec = ""
	}
	return element, property, fmtspec
}

func getField(header string, begin, end int) string {
	if begin < end {
		return strings.TrimSpace(header[begin:end])
	} else {
		return ""
	}
}
