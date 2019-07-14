// Package file implements a file writer.
package file

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	gos "github.com/gratonos/goutil/os"
	"github.com/gratonos/gxlog/iface"
)

const (
	checkInterval = time.Second * 5
	dateFormat    = "%04d%02d%02d"
	timeFormat    = "%02d%02d%02d.%06d"
	extension     = ".log"
	dirPerm       = 0770
)

type Writer struct {
	dir         string
	maxFileSize int64

	writer    io.WriteCloser
	path      string
	checkTime time.Time
	yearDay   int
	fileSize  int64
	lock      sync.Mutex
}

func Open(config Config) (*Writer, error) {
	config.SetDefaults()

	if err := checkDir(config.Dir); err != nil {
		return nil, fmt.Errorf("writer/file.Open: %v", err)
	}

	return &Writer{
		dir:         config.Dir,
		maxFileSize: config.MaxFileSize,
	}, nil
}

func (this *Writer) Close() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.closeFile(); err != nil {
		return fmt.Errorf("writer/file.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *iface.Record) error {
	this.lock.Lock()

	err := this.checkFile(record)
	if err == nil {
		var n int
		n, err = this.writer.Write(bs)
		this.fileSize += int64(n)
	}

	this.lock.Unlock()

	if err != nil {
		return fmt.Errorf("writer/file.Write: %v", err)
	}
	return nil
}

func (this *Writer) Dir() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.dir
}

func (this *Writer) SetDir(dir string) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if dir == this.dir {
		return nil
	}
	if err := checkDir(dir); err != nil {
		return fmt.Errorf("writer/file.SetDir: %v", err)
	}
	if err := this.closeFile(); err != nil {
		return fmt.Errorf("writer/file.SetDir: %v", err)
	}
	this.dir = dir
	return nil
}

func (this *Writer) MaxFileSize() int64 {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.maxFileSize
}

func (this *Writer) SetMaxFileSize(size int64) error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if size <= 0 {
		return errors.New("writer/file.SetMaxFileSize: size must be positive")
	}
	this.maxFileSize = size
	return nil
}

func (this *Writer) checkFile(record *iface.Record) error {
	if this.writer == nil ||
		this.yearDay != record.Time.YearDay() ||
		this.fileSize >= this.maxFileSize {
		return this.createFile(record)
	} else if time.Since(this.checkTime) >= checkInterval {
		this.checkTime = time.Now()
		ok, err := gos.FileExists(this.path)
		if err != nil {
			return err
		}
		if !ok {
			return this.createFile(record)
		}
	}
	return nil
}

func (this *Writer) createFile(record *iface.Record) error {
	if err := this.closeFile(); err != nil {
		return err
	}

	date := fmt.Sprintf(dateFormat, record.Time.Year(), record.Time.Month(), record.Time.Day())
	dir := filepath.Join(this.dir, date)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return err
	}

	clock := fmt.Sprintf(timeFormat, record.Time.Hour(), record.Time.Minute(),
		record.Time.Second(), record.Time.Nanosecond()/1000)
	filename := clock + extension
	path := filepath.Join(dir, filename)
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	this.writer = file
	this.path = path
	this.yearDay = record.Time.YearDay()
	this.fileSize = 0

	return nil
}

func (this *Writer) closeFile() error {
	if this.writer != nil {
		if err := this.writer.Close(); err != nil {
			return err
		}
		this.writer = nil
	}
	return nil
}

func checkDir(dir string) error {
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return err
	}
	return syscall.Access(dir, 7 /* R_OK | W_OK | X_OK */)
}
