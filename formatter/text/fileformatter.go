package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/formatter/internal/util"
	"github.com/gratonos/gxlog/iface"
)

type fileFormatter struct {
	segments int
	fmtspec  string
}

func newFileFormatter(property, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &fileFormatter{
		segments: segments,
		fmtspec:  fmtspec,
	}
}

func (formatter *fileFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	file := util.LastSegments(record.File, formatter.segments, '/')
	if formatter.fmtspec == "%s" {
		return append(buf, file...)
	}
	return append(buf, fmt.Sprintf(formatter.fmtspec, file)...)
}
