package text_test

import (
	"fmt"
	"path"
	"path/filepath"
	"testing"
	"time"

	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/iface"
)

const (
	tmplLayout      = "2006-01-02 15:04:05.000000000"
	tmplDate        = "2018-08-01"
	tmplTime        = "07:12:07"
	tmplDecimal     = "235605270"
	tmplLevel       = iface.Info
	tmplFile        = "/home/test/data/src/go/workspace/src/github.com/gratonos/gxlog/logger.go"
	tmplLine        = 64
	tmplPkg         = "github.com/gratonos/gxlog"
	tmplFunc        = "Test"
	tmplMsg         = "testing"
	tmplPrefix      = "**** "
	tmplContextPair = "(k1: v1) (k2: v2)"
	tmplContextList = "k1: v1, k2: v2"
)

var tmplTimestamp time.Time

func init() {
	var err error
	tmplTimestamp, err = time.ParseInLocation(tmplLayout,
		tmplDate+" "+tmplTime+"."+tmplDecimal, time.Local)
	if err != nil {
		panic(err)
	}
}

func TestFullHeader(t *testing.T) {
	formatter := text.New(text.Config{Header: text.FullHeader})
	expect := fmt.Sprintf("%s %s.%s %s %s:%d %s.%s %s[%s] %s\n",
		tmplDate, tmplTime, tmplDecimal[:6], "INFO ", tmplFile, tmplLine,
		tmplPkg, tmplFunc, tmplPrefix, tmplContextPair, tmplMsg)
	testFormat(t, formatter, tmplRecord(), expect)
}

func TestCompactHeader(t *testing.T) {
	formatter := text.New(text.Config{Header: text.CompactHeader})
	expect := fmt.Sprintf("%s.%s %s %s:%d %s.%s %s[%s] %s\n",
		tmplTime, tmplDecimal[:6], "INFO ", filepath.Base(tmplFile), tmplLine,
		tmplPkg, tmplFunc, tmplPrefix, tmplContextPair, tmplMsg)
	testFormat(t, formatter, tmplRecord(), expect)
}

func TestSyslogHeader(t *testing.T) {
	formatter := text.New(text.Config{Header: text.SyslogHeader})
	expect := fmt.Sprintf("%s:%d %s.%s %s[%s] %s\n",
		filepath.Base(tmplFile), tmplLine, tmplPkg, tmplFunc, tmplPrefix,
		tmplContextPair, tmplMsg)
	testFormat(t, formatter, tmplRecord(), expect)
}

func TestCustomHeader(t *testing.T) {
	header := "{{time:time.ns}} {{level:char}} {{file}}:{{line%05d}} " +
		"{{pkg:1}}.{{func}} {{prefix}}[{{context:list}}] {{msg%20s}}\n"
	formatter := text.New(text.Config{Header: header})
	expect := fmt.Sprintf("%s.%s %s %s:%05d %s.%s %s[%s] %20s\n",
		tmplTime, tmplDecimal, "I", tmplFile, tmplLine, path.Base(tmplPkg),
		tmplFunc, tmplPrefix, tmplContextList, tmplMsg)
	testFormat(t, formatter, tmplRecord(), expect)
}

func TestBizarreHeader(t *testing.T) {
	header := "xx{{unknown}} {static} {{unknown}} {{level : char}} {{pkg|1}} " +
		"[{{context:dot}}] {{ msg %20s }}yy"
	formatter := text.New(text.Config{Header: header})
	expect := fmt.Sprintf("xx {static}  %s  [%s] %20syy", "I", tmplContextPair, tmplMsg)
	testFormat(t, formatter, tmplRecord(), expect)
}

func TestColor(t *testing.T) {
	formatter := text.New(text.Config{Header: "{{msg}}", Coloring: true})
	expect := fmt.Sprintf("\033[%dm%s\033[0m", text.Magenta, tmplMsg)
	testFormat(t, formatter, tmplRecord(), expect)

	record := tmplRecord()
	record.Level = iface.Warn
	record.Mark = false
	formatter.MapColors(map[iface.Level]text.Color{
		iface.Warn: text.Blue,
	})
	expect = fmt.Sprintf("\033[%dm%s\033[0m", text.Blue, tmplMsg)
	testFormat(t, formatter, record, expect)

	record.Level = iface.Error
	formatter.SetColor(iface.Error, text.Yellow)
	expect = fmt.Sprintf("\033[%dm%s\033[0m", text.Yellow, tmplMsg)
	testFormat(t, formatter, record, expect)
}

func testFormat(t *testing.T, formatter iface.Formatter, record *iface.Record, expect string) {
	output := string(formatter.Format(record))
	if output != expect {
		t.Errorf("testFormat:\noutput: %q\nexpect: %q", output, expect)
	}
}

func tmplRecord() *iface.Record {
	return &iface.Record{
		Time:   tmplTimestamp,
		Level:  tmplLevel,
		File:   tmplFile,
		Line:   tmplLine,
		Pkg:    tmplPkg,
		Func:   tmplFunc,
		Msg:    tmplMsg,
		Prefix: tmplPrefix,
		Contexts: []iface.Context{
			{Key: "k1", Value: "v1"},
			{Key: "k2", Value: "v2"},
		},
		Mark: true,
	}
}
