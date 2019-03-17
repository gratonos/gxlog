package main

import (
	"strings"

	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/logger"
)

var log = gxlog.Logger()

func main() {
	testFilterLogic()
}

func testFilterLogic() {
	log.SetFilter(logger.Or(important, logger.And(useful, interesting)))
	log.Error("error") // this will be output
	log.Warn("warn")
	log.Trace("trace, funny")
	log.Info("info, funny") // this will be output
}

func important(record *iface.Record) bool {
	return record.Level >= iface.Error
}

func useful(record *iface.Record) bool {
	return record.Level >= iface.Info
}

func interesting(record *iface.Record) bool {
	return strings.Contains(record.Msg, "funny")
}
