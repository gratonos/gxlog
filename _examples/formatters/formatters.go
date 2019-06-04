package main

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/formatter"
	"github.com/gratonos/gxlog/formatter/json"
	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
	testCustomFormatter()
	testTextFormatter()
	testJSONFormatter()
}

func testCustomFormatter() {
	fn := formatter.Func(func(record *iface.Record) []byte {
		return []byte(record.Msg + "\n")
	})
	log.SetSlotFormatter(logger.Slot0, fn)
	log.Info("a simple formatter that just returns the msg of a record")
}

func testTextFormatter() {
	textFmt := gxlog.Formatter()
	log.SetSlotFormatter(logger.Slot0, textFmt)

	log.Trace("green")
	log.Warn("yellow")
	log.Error("red")
	log.Mark(true).Error("magenta")

	textFmt.SetHeader(text.FullHeader)
	textFmt.SetColor(iface.Trace, text.Blue)
	textFmt.MapColors(map[iface.Level]text.Color{
		iface.Warn:  text.Red,
		iface.Error: text.Magenta,
	})
	textFmt.SetMarkColor(text.White)
	log.Trace("blue")
	log.Warn("red")
	log.Error("magenta")
	log.Mark(true).Error("white")

	header := "{{time:time}} {{level:char}} {{file:2%q}}:{{line:%05d}} {{msg:%20s}}\n"
	textFmt.SetHeader(header)
	textFmt.SetColoring(false)
	log.Trace("default color")
}

func testJSONFormatter() {
	jsonFmt := json.New(json.Config{FileSegs: 1})
	log.SetSlotFormatter(logger.Slot0, jsonFmt)
	log.Trace("json")

	jsonFmt.SetFileSegs(0)
	log.Trace("json updated")
}
