package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/iface"
)

type lineFormatter struct {
	fmtspec string
}

func newLineFormatter(_, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%d"
	}
	return &lineFormatter{fmtspec: fmtspec}
}

func (formatter *lineFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if formatter.fmtspec == "%d" {
		return strconv.AppendInt(buf, int64(record.Line), 10)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, record.Line)...)
}
