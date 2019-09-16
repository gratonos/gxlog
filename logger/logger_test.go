package logger_test

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/logger"
)

type externalChecker func(log string) error

const (
	goroutines = 8
	logCount   = 10000
	msgLen     = 64
)

const (
	dateRegexp     = "[0-9]{4}-[0-9]{2}-[0-9]{2}"          // e.g. 2019-09-15
	timeRegexp     = "[0-9]{2}:[0-9]{2}:[0-9]{2}.[0-9]{6}" // e.g. 11:05:33.722028
	contextsRegexp = "\\[\\(id: ([0-9])\\)\\]"             // e.g. [(id: 1)]
	msgRegexp      = "([0-9]+)"                            // e.g. 1111...1
)

var logRegexp *regexp.Regexp

func init() {
	// e.g. 2019-09-15 11:05:33.722028 [(id: 1)] 1111...1
	fullRegexp := fmt.Sprintf("^%s %s %s %s$", dateRegexp, timeRegexp, contextsRegexp, msgRegexp)
	logRegexp = regexp.MustCompile(fullRegexp)
}

func TestConcurrency(t *testing.T) {
	writer := newTestWriter(newChecker())
	log := newLogger(writer, t)

	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)

		go func(id int) {
			log := log.WithStatics("id", id)
			msg := strings.Repeat(strconv.Itoa(id), msgLen)

			for n := 0; n < logCount; n++ {
				log.Info(msg)
			}

			wg.Done()
		}(i)
	}

	wg.Wait()

	expectLogs := goroutines * logCount
	gotLogs := writer.TotalLogs()
	if expectLogs != gotLogs {
		t.Errorf("log missing, expect %d, got %d", expectLogs, gotLogs)
	}
}

func newChecker() externalChecker {
	return func(log string) error {
		indexes := logRegexp.FindStringSubmatchIndex(log)
		if indexes == nil {
			return errors.New("regexp unmatch")
		}
		// 2019-09-15 11:05:33.722028 [(id: 1)] 1111...1
		// 0<------------------------------------------>1
		//                                  23
		//                                      4<----->5
		if len(indexes) != 6 {
			panic("corrupted regexp")
		}
		for _, index := range indexes {
			if index < 0 {
				panic("sub-regexp unmatch")
			}
		}
		id := log[indexes[2]:indexes[3]]
		msg := log[indexes[4]:indexes[5]]

		if strings.Repeat(id, msgLen) != msg {
			return errors.New("msg corrupted")
		}
		return nil
	}
}

func newLogger(writer iface.Writer, t *testing.T) *logger.Logger {
	formatter := text.New(text.Config{
		Header: "{{time}} [{{context}}] {{msg}}",
	})

	handler := func(_ []byte, _ *iface.Record, err error) {
		t.Errorf("%v", err)
	}

	log := logger.New(logger.Config{})
	log.SetSlot(logger.Slot0, logger.Slot{
		Formatter:    formatter,
		Writer:       writer,
		ErrorHandler: handler,
	})
	return log
}
