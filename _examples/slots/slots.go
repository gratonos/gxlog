package main

import (
	"fmt"

	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/formatter"
	"github.com/gratonos/gxlog/formatter/json"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
	testSlots()
	testSlotsLevel()
}

func testSlots() {
	log.Info("this will be printed once")

	log.CopySlot(logger.Slot1, logger.Slot0)
	log.Info("this will be printed twice")

	log.SetSlotFormatter(logger.Slot1, json.New(json.Config{}))
	log.Info("this will be printed in text format and json format")

	log.SwapSlot(logger.Slot0, logger.Slot1)
	log.Info("json first and then text")

	log.MoveSlot(logger.Slot0, logger.Slot1)
}

func testSlotsLevel() {
	log.SetSlotLevel(logger.Slot0, iface.Warn)
	log.Info("this will not be printed")
	log.Warn("this will be printed")

	log.SetSlotLevel(logger.Slot0, iface.Trace)
	// ATTENTION: Do NOT call any method of the Logger in a Formatter, Writer
	// or Filter, or it may deadlock.
	hook := formatter.Func(func(record *iface.Record) []byte {
		fmt.Println("hooks:", record.Msg)
		return nil
	})
	filter := func(record *iface.Record) bool {
		return record.Aux.Mark
	}
	slot := logger.Slot{
		Formatter:    hook,
		Level:        iface.Warn,
		Filter:       filter,
		ErrorHandler: logger.Report,
	}
	log.SetSlot(logger.Slot0, slot)

	log.Mark(true).Info("marked, but info")
	log.Error("error, but not marked")
	log.Mark(true).Warn("warn and marked")
}
