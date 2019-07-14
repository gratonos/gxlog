package text

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

type prefixFormatter struct {
	fmtspec string
}

func newPrefixFormatter(_, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &prefixFormatter{fmtspec: fmtspec}
}

func (this *prefixFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if this.fmtspec == "%s" {
		return append(buf, record.Prefix...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, record.Prefix)...)
	}
}
