package main

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
	testAuxiliary()
	testDynamicContext()
}

func testAuxiliary() {
	log.Prefix("**** ").Mark(true).Contexts("k1", "v1", "k2", "v2").
		Info("prefix, mark and contexts")
	log.Info("no prefix, mark or contexts")

	func() {
		log := log.Contexts("k3", "v3")
		log.Info("outer enter")
		func() {
			log := log.Contexts("k4", "v4")
			log.Info("inner")
		}()
		log.Info("outer leave")
	}()
}

func testDynamicContext() {
	// ATTENTION: You SHOULD be very careful about concurrency safety or deadlocks
	// with dynamic contexts.
	n := 0
	fn := logger.Dynamic(func(interface{}) interface{} {
		// Do NOT call any method of the Logger in the function, or it may deadlock.
		n++
		return n
	})
	clog := log.Contexts("static", n, "dynamic", fn)
	clog.Info("dynamic one")
	clog.Info("dynamic two")
}
