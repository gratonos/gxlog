package text

import (
	"fmt"
	"strconv"

	"github.com/gratonos/gxlog/formatter/internal/util"
	"github.com/gratonos/gxlog/iface"
)

type pkgFormatter struct {
	segments int
	fmtspec  string
}

func newPkgFormatter(property, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	segments, _ := strconv.Atoi(property)
	return &pkgFormatter{
		segments: segments,
		fmtspec:  fmtspec,
	}
}

func (this *pkgFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	pkg := util.LastSegments(record.Pkg, this.segments, '/')
	if this.fmtspec == "%s" {
		return append(buf, pkg...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, pkg)...)
	}
}
