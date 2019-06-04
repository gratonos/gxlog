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

func (formatter *Formatter) Header() string {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.header
}

func (formatter *Formatter) SetHeader(header string) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.header = header
	formatter.appenders = formatter.appenders[:0]
	var staticText string
	for header != "" {
		indexes := headerRegexp.FindStringSubmatchIndex(header)
		if indexes == nil {
			break
		}
		begin, end := indexes[0], indexes[1]
		staticText += header[:begin]
		element, property, fmtspec := extractElement(indexes, header)
		if formatter.addAppender(element, property, fmtspec, staticText) {
			staticText = ""
		}
		header = header[end:]
	}
	formatter.suffix = staticText + header
}

func (formatter *Formatter) Coloring() bool {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.coloring
}

func (formatter *Formatter) SetColoring(ok bool) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.coloring = ok
}

func (formatter *Formatter) Color(level iface.Level) Color {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.colorMgr.Color(level)
}

func (formatter *Formatter) SetColor(level iface.Level, color Color) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.colorMgr.SetColor(level, color)
}

func (formatter *Formatter) MapColors(colorMap map[iface.Level]Color) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.colorMgr.MapColors(colorMap)
}

func (formatter *Formatter) MarkColor() Color {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	return formatter.colorMgr.MarkColor()
}

func (formatter *Formatter) SetMarkColor(color Color) {
	formatter.lock.Lock()
	defer formatter.lock.Unlock()

	formatter.colorMgr.SetMarkColor(color)
}

func (formatter *Formatter) Format(record *iface.Record) []byte {
	formatter.lock.Lock()

	var left, right []byte
	if formatter.coloring {
		if record.Aux.Mark {
			left, right = formatter.colorMgr.MarkColorEars()
		} else {
			left, right = formatter.colorMgr.ColorEars(record.Level)
		}
	}

	buf := formatter.buf[:0]
	buf = append(buf, left...)
	for _, appender := range formatter.appenders {
		buf = appender.AppendHeader(buf, record)
	}
	buf = append(buf, formatter.suffix...)
	buf = append(buf, right...)
	formatter.buf = buf

	formatter.lock.Unlock()
	return buf
}

func (formatter *Formatter) addAppender(element, property, fmtspec, staticText string) bool {
	appender := newHeaderAppender(element, property, fmtspec, staticText)
	if appender == nil {
		return false
	}
	formatter.appenders = append(formatter.appenders, appender)
	return true
}

func extractElement(indexes []int, header string) (element, property, fmtspec string) {
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
	}
	return ""
}
