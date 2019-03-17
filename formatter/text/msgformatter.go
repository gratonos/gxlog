package text

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

type msgFormatter struct {
	fmtspec string
}

func newMsgFormatter(_, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &msgFormatter{fmtspec: fmtspec}
}

func (formatter *msgFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if formatter.fmtspec == "%s" {
		return append(buf, record.Msg...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, record.Msg)...)
}
