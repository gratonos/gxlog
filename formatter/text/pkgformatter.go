package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/formatter/internal/util"
	"github.com/gratonos/gxlog/iface"
)

type pkgFormatter struct {
	segments int
	fmtstr  string
}

func newPkgFormatter(property, fmtstr string) elementFormatter {
	if fmtstr == "" {
		fmtstr = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &pkgFormatter{
		segments: segments,
		fmtstr:  fmtstr,
	}
}

func (this *pkgFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	pkg := util.LastSegments(record.Pkg, this.segments, '/')
	if this.fmtstr == "%s" {
		return append(buf, pkg...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtstr, pkg)...)
	}
}
