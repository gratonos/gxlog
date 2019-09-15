package logger_test

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/logger"
	"github.com/gratonos/gxlog/writer"
)

const msg = "0123456789012345678901234567890123456789012345678901234567890123"

func BenchmarkLogger(b *testing.B) {
	log := logger.New(logger.Config{})
	log.SetSlot(logger.Slot0, logger.Slot{
		Formatter: text.New(text.Config{
			Header: "{{time}} {{file:1}}:{{line}}: {{msg}}\n",
		}),
		Writer: writer.Wrap(ioutil.Discard),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Info(msg)
	}
}

func BenchmarkStdLogger(b *testing.B) {
	log := log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println(msg)
	}
}
