package text

import (
	"fmt"
	"strings"

	"github.com/gratonos/gxlog/iface"
)

var levelNames = []string{
	iface.Trace: "TRACE",
	iface.Debug: "DEBUG",
	iface.Info:  "INFO ",
	iface.Warn:  "WARN ",
	iface.Error: "ERROR",
	iface.Fatal: "FATAL",
}

var levelCharNames = []string{
	iface.Trace: "T",
	iface.Debug: "D",
	iface.Info:  "I",
	iface.Warn:  "W",
	iface.Error: "E",
	iface.Fatal: "F",
}

type levelFormatter struct {
	nameList []string
	fmtstr   string
}

func newLevelFormatter(property, fmtstr string) elementFormatter {
	if fmtstr == "" {
		fmtstr = "%s"
	}
	return &levelFormatter{
		nameList: selectNameList(property),
		fmtstr:   fmtstr,
	}
}

func (this *levelFormatter) FormatElement(buf []byte, record *iface.Record) []byte {
	name := this.nameList[record.Level]
	if this.fmtstr == "%s" {
		return append(buf, name...)
	} else {
		return append(buf, fmt.Sprintf(this.fmtstr, name)...)
	}
}

func selectNameList(property string) []string {
	if strings.ToLower(property) == "char" {
		return levelCharNames
	} else {
		return levelNames
	}
}
