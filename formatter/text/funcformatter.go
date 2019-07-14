package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/formatter/internal/util"
	"github.com/gratonos/gxlog/iface"
)

type funcFormatter struct {
	segments int
	fmtspec  string
}

func newFuncFormatter(property, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &funcFormatter{
		segments: segments,
		fmtspec:  fmtspec,
	}
}

func (this *funcFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	fn := util.LastSegments(record.Func, this.segments, '.')
	if this.fmtspec == "%s" {
		return append(buf, fn...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, fn)...)
	}
}
