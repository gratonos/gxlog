package text

import (
	"fmt"
	"strings"

	"github.com/gratonos/gxlog/iface"
)

type contextFormatter struct {
	formatter func([]byte, []iface.Context) []byte
	fmtspec   string
	buf       []byte
}

func newContextFormatter(property, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &contextFormatter{
		formatter: selectFormatter(property),
		fmtspec:   fmtspec,
	}
}

func (this *contextFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if this.fmtspec == "%s" {
		return this.formatter(buf, record.Contexts)
	} else {
		this.buf = this.buf[:0]
		this.buf = this.formatter(this.buf, record.Contexts)
		return append(buf, fmt.Sprintf(this.fmtspec, this.buf)...)
	}
}

func selectFormatter(property string) func([]byte, []iface.Context) []byte {
	if strings.ToLower(property) == "list" {
		return formatList
	} else {
		return formatPair
	}
}

func formatPair(buf []byte, contexts []iface.Context) []byte {
	left := "("
	for _, ctx := range contexts {
		buf = append(buf, left...)
		buf = append(buf, ctx.Key...)
		buf = append(buf, ": "...)
		buf = append(buf, ctx.Value...)
		buf = append(buf, ')')
		left = " ("
	}
	return buf
}

func formatList(buf []byte, contexts []iface.Context) []byte {
	begin := ""
	for _, ctx := range contexts {
		buf = append(buf, begin...)
		buf = append(buf, ctx.Key...)
		buf = append(buf, ": "...)
		buf = append(buf, ctx.Value...)
		begin = ", "
	}
	return buf
}
