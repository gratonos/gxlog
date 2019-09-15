package logger_test

import (
	"fmt"

	"github.com/gratonos/gxlog/iface"
)

type testWriter struct {
	checker   externalChecker
	prevLog   string
	totalLogs int
}

func newTestWriter(checker externalChecker) *testWriter {
	return &testWriter{
		checker: checker,
	}
}

func (this *testWriter) Write(bs []byte, _ *iface.Record) error {
	log := string(bs)

	if err := this.checker(log); err != nil {
		return fmt.Errorf("%v, log: %s", err, log)
	}
	if this.prevLog >= log {
		return fmt.Errorf("out of order: prev: %s, cur: %s", this.prevLog, log)
	}

	this.prevLog = log
	this.totalLogs++
	return nil
}

func (this *testWriter) TotalLogs() int {
	return this.totalLogs
}
