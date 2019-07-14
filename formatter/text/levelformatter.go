package text

import (
	"fmt"
	"strings"

	"github.com/gratonos/gxlog/iface"
)

var levelDesc = []string{
	iface.Trace: "TRACE",
	iface.Debug: "DEBUG",
	iface.Info:  "INFO ",
	iface.Warn:  "WARN ",
	iface.Error: "ERROR",
	iface.Fatal: "FATAL",
}

var levelDescChar = []string{
	iface.Trace: "T",
	iface.Debug: "D",
	iface.Info:  "I",
	iface.Warn:  "W",
	iface.Error: "E",
	iface.Fatal: "F",
}

type levelFormatter struct {
	descList []string
	fmtspec  string
}

func newLevelFormatter(property, fmtspec string) elementFormatter {
	if fmtspec == "" {
		fmtspec = "%s"
	}
	return &levelFormatter{
		descList: selectDescList(property),
		fmtspec:  fmtspec,
	}
}

func (this *levelFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	desc := this.descList[record.Level]
	if this.fmtspec == "%s" {
		return append(buf, desc...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtspec, desc)...)
	}
}

func selectDescList(property string) []string {
	if strings.ToLower(property) == "char" {
		return levelDescChar
	} else {
		return levelDesc
	}
}
