package text

import (
	"github.com/gratonos/gxlog/iface"
)

type elementFormatter interface {
	FormatElement(buf []byte, record *iface.Record) []byte
}

var newFormatterFuncMap = map[string]func(property, fmtstr string) elementFormatter{
	"time":    newTimeFormatter,
	"level":   newLevelFormatter,
	"file":    newFileFormatter,
	"line":    newLineFormatter,
	"pkg":     newPkgFormatter,
	"func":    newFuncFormatter,
	"msg":     newMsgFormatter,
	"prefix":  newPrefixFormatter,
	"context": newContextFormatter,
}

type headerAppender struct {
	formatter  elementFormatter
	staticText string
}

func newHeaderAppender(element, property, fmtstr, staticText string) *headerAppender {
	newFunc := newFormatterFuncMap[element]
	if newFunc == nil {
		return nil
	}
	return &headerAppender{
		formatter:  newFunc(property, fmtstr),
		staticText: staticText,
	}
}

func (this *headerAppender) AppendHeader(buf []byte, record *iface.Record) []byte {
	buf = append(buf, this.staticText...)
	return this.formatter.FormatElement(buf, record)
}
