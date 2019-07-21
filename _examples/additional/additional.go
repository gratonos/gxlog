package main

import (
	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
	testAdditional()
	testDynamicContext()
}

func testAdditional() {
	log.WithPrefix("**** ").WithMark(true).WithStatics("k1", "v1", "k2", "v2").
		Info("prefix, mark and contexts")
	log.Info("no prefix, mark or contexts")

	func() {
		log := log.WithStatics("k3", "v3")
		log.Info("outer enter")
		func() {
			log := log.WithStatics("k4", "v4")
			log.Info("inner")
		}()
		log.Info("outer leave")
	}()
}

func testDynamicContext() {
	// ATTENTION: You SHOULD be very careful about concurrency safety or deadlocks
	// with dynamic contexts.
	n := 0
	fn := logger.Dynamic(func() interface{} {
		// Do NOT call any method of the Logger in the function, or it may deadlock.
		n++
		return n
	})
	clog := log.WithStatics("static", n).WithDynamics("dynamic", fn)
	clog.Info("dynamic one")
	clog.Info("dynamic two")
}
