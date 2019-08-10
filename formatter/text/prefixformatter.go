package text

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

type prefixFormatter struct {
	fmtstr string
}

func newPrefixFormatter(_, fmtstr string) elementFormatter {
	if fmtstr == "" {
		fmtstr = "%s"
	}
	return &prefixFormatter{fmtstr: fmtstr}
}

func (this *prefixFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if this.fmtstr == "%s" {
		return append(buf, record.Prefix...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtstr, record.Prefix)...)
	}
}
