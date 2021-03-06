package main

import (
	"fmt"
	"os"

	"github.com/gratonos/gxlog"
	"github.com/gratonos/gxlog/formatter/text"
	"github.com/gratonos/gxlog/iface"
	"github.com/gratonos/gxlog/logger"
	"github.com/gratonos/gxlog/writer"
	"github.com/gratonos/gxlog/writer/file"
	"github.com/gratonos/gxlog/writer/syslog"
	"github.com/gratonos/gxlog/writer/usock"
)

var log = gxlog.Logger()

func main() {
	testWrappers()
	testUSockWriter()

	gxlog.Formatter().SetColoring(false)

	testFileWriter()
	testSyslogWriter()
}

func testWrappers() {
	fn := writer.Func(func(bs []byte, _ *iface.Record) error {
		_, err := os.Stderr.Write(bs)
		return err
	})
	log.SetSlotWriter(logger.Slot0, fn)
	log.Info("a simple writer that just writes to os.Stderr")

	// another equivalent way
	log.SetSlotWriter(logger.Slot0, writer.Wrap(os.Stderr))
	log.Info("writer wrapper of os.Stderr")
}

func testUSockWriter() {
	// Shell expansion is NOT supported. Thus, ~, $var and so on will NOT be expanded.
	writer, err := usock.Open("/tmp/gxlog/usock")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer writer.Close()

	log.SetSlotWriter(logger.Slot0, writer)

	// Use "netcat -U /tmp/gxlog/usock" to watch logs.
	// for i := 0; i < 1024; i++ {
	// 	log.Info(i)
	// 	time.Sleep(time.Second)
	// }
}

func testFileWriter() {
	// Shell expansion is NOT supported. Thus, ~, $var and so on will NOT be expanded.
	writer, err := file.Open(file.Config{Dir: "/tmp/gxlog"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer writer.Close()

	log.SetSlotWriter(logger.Slot0, writer)
	log.Info("this will be output to a file")
}

func testSyslogWriter() {
	gxlog.Formatter().SetHeader(text.SyslogHeader)

	writer, err := syslog.Open(syslog.Config{Tag: "gxlog", Facility: syslog.FacUser})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer writer.Close()

	log.SetSlotWriter(logger.Slot0, writer)
	log.Info("this will be output to syslog")

	writer.MapSeverities(map[iface.Level]syslog.Severity{
		iface.Info: syslog.SevErr,
	})
	log.Info("this will be severity err")
}
