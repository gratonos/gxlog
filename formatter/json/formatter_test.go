package json_test

import (
	"bytes"
	ejson "encoding/json"
	"testing"
	"time"

	"github.com/gratonos/gxlog/formatter/json"
	"github.com/gratonos/gxlog/iface"
)

const (
	tmplLayout      = "2006-01-02 15:04:05.000000000"
	tmplDate        = "2018-08-01"
	tmplTime        = "07:12:07"
	tmplFraction    = "235605270"
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
		tmplDate+" "+tmplTime+"."+tmplFraction, time.Local)
	if err != nil {
		panic(err)
	}
}

func TestFull(t *testing.T) {
	formatter := json.New(json.Config{})
	expect := jsonMarshal(tmplRecord())
	output := formatter.Format(tmplRecord())
	if !bytes.Equal(expect, output) {
		t.Errorf("TestFull:\noutput: %q\nexpect: %q", output, expect)
	}
}

func jsonMarshal(record *iface.Record) []byte {
	bs, err := ejson.Marshal(record)
	if err != nil {
		panic("json Marshal failed")
	}
	return append(bs, '\n')
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
