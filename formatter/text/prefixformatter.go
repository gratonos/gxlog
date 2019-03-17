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

func (formatter *prefixFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if formatter.fmtspec == "%s" {
		return append(buf, record.Aux.Prefix...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, record.Aux.Prefix)...)
}
