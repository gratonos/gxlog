package text

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

type msgFormatter struct {
	fmtstr string
}

func newMsgFormatter(_, fmtstr string) elementFormatter {
	if fmtstr == "" {
		fmtstr = "%s"
	}
	return &msgFormatter{fmtstr: fmtstr}
}

func (this *msgFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	if this.fmtstr == "%s" {
		return append(buf, record.Msg...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtstr, record.Msg)...)
	}
}
