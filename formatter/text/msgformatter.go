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

func (this *msgFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if this.fmtspec == "%s" {
		return append(buf, record.Msg...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, record.Msg)...)
	}
}
