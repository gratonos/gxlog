package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/formatter/internal/util"
	"github.com/gratonos/gxlog/iface"
)

type fileFormatter struct {
	segments int
	fmtstr  string
}

func newFileFormatter(property, fmtstr string) elementFormatter {
	if fmtstr == "" {
		fmtstr = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &fileFormatter{
		segments: segments,
		fmtstr:  fmtstr,
	}
}

func (this *fileFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	file := util.LastSegments(record.File, this.segments, '/')
	if this.fmtstr == "%s" {
		return append(buf, file...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtstr, file)...)
	}
}
