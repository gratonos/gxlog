package logger

import (
	"fmt"
	"log"

	"github.com/gratonos/gxlog/iface"
)

// Do NOT call any method of the Logger within a ErrorHandler, or it may deadlock.
type ErrorHandler func(bs []byte, record *iface.Record, err error)

const outputDepthOffset = 6

func Report(_ []byte, _ *iface.Record, err error) {
	_ = log.Output(outputDepthOffset, fmt.Sprintln("gxlog error:", err))
}

func ReportDetails(bs []byte, _ *iface.Record, err error) {
	_ = log.Output(outputDepthOffset, fmt.Sprintf("gxlog error: %v, log: %s", err, bs))
}
