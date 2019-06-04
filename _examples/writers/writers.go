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
	"github.com/gratonos/gxlog/writer/unix"
)

var log = gxlog.Logger()

func main() {
	testWrappers()
	testUnixWriter()

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

func testUnixWriter() {
	// Shell expansion is NOT supported. Thus, ~, $var and so on will NOT be expanded.
	wt, err := unix.Open("/tmp/gxlog/unixdomain")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer wt.Close()

	log.SetSlotWriter(logger.Slot0, wt)

	// Use "netcat -U /tmp/gxlog/unixdomain" to watch logs.
	// for i := 0; i < 1024; i++ {
	// 	log.Info(i)
	// 	time.Sleep(time.Second)
	// }
}

func testFileWriter() {
	// Shell expansion is NOT supported. Thus, ~, $var and so on will NOT be expanded.
	wt, err := file.Open(file.Config{Path: "/tmp/gxlog"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer wt.Close()

	log.SetSlotWriter(logger.Slot0, wt)
	log.Info("this will be output to a file")
}

func testSyslogWriter() {
	gxlog.Formatter().SetHeader(text.SyslogHeader)

	wt, err := syslog.Open(syslog.Config{Tag: "gxlog", Facility: syslog.FacUser})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer wt.Close()

	log.SetSlotWriter(logger.Slot0, wt)
	log.Info("this will be output to syslog")

	wt.MapSeverities(map[iface.Level]syslog.Severity{
		iface.Info: syslog.SevErr,
	})
	log.Info("this will be severity err")
}
