package main

import (
	"time"

	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/iface"
)

var log = gxlog.Logger()

func main() {
	testLevel()
	// testPanic()
	testTime()
	testLog()
	testEError()
	testEErrorf()
}

func testLevel() {
	log.Trace("test Trace")
	log.Tracef("%s", "test Tracef")
	log.Debug("test Debug")
	log.Debugf("%s", "test Debugf")
	log.Info("test Info")
	log.Infof("%s", "test Infof")
	log.Warn("test Warn")
	log.Warnf("%s", "test Warnf")
	log.Error("test Error")
	log.Errorf("%s", "test Errorf")
	// log.Fatal("test Fatal")
	// log.Fatalf("%s", "test Fatalf")
}

func testPanic() {
	log.Panic("test Panic")
	log.Panicf("%s", "test Panicf")
}

func testTime() {
	done := log.Timing(iface.Trace, "test Time")
	time.Sleep(200 * time.Millisecond)
	done()

	defer log.Timingf(iface.Trace, "%s", "test Timef")()
	time.Sleep(400 * time.Millisecond)
}

func testLog() {
	log.Log(0, iface.Info, "test Log")
	log.Logf(1, iface.Warn, "%s: %d", "test Logf", 1)
	log.Logf(-1, iface.Warn, "%s: %d", "test Logf", -1)
}

func testEError() error {
	return log.EError("an error")
}

func testEErrorf() error {
	return log.EErrorf("%s", "another error")
}
