package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/iface"
)

type lineFormatter struct {
	fmtstr string
}

func newLineFormatter(_, fmtstr string) elementFormatter {
	if fmtstr == "" {
		fmtstr = "%d"
	}
	return &lineFormatter{fmtstr: fmtstr}
}

func (this *lineFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if this.fmtstr == "%d" {
		return strconv.AppendInt(buf, int64(record.Line), 10)
	} else {
		return append(buf, fmt.Sprintf(this.fmtstr, record.Line)...)
	}
}
