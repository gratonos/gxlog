package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/formatter/internal/util"
	"github.com/gratonos/gxlog/iface"
)

type funcFormatter struct {
	segments int
	fmtstr  string
}

func newFuncFormatter(property, fmtstr string) elementFormatter {
	if fmtstr == "" {
		fmtstr = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &funcFormatter{
		segments: segments,
		fmtstr:  fmtstr,
	}
}

func (this *funcFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	fn := util.LastSegments(record.Func, this.segments, '.')
	if this.fmtstr == "%s" {
		return append(buf, fn...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtstr, fn)...)
	}
}
