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
	path        string
	maxFileSize int64

	writer    io.WriteCloser
	pathname  string
	checkTime time.Time
	day       int
	fileSize  int64
	lock      sync.Mutex
}

func Open(path string, maxFileSize int64) (*Writer, error) {
	if path == "" {
		return nil, errors.New("writer/file.Open: path must not be empty")
	}
	if maxFileSize <= 0 {
		return nil, errors.New("writer/file.Open: maxFileSize must be positive")
	}
	if err := checkPath(path); err != nil {
		return nil, fmt.Errorf("writer/file.Open: %v", err)
	}
	return &Writer{
		path:        path,
		maxFileSize: maxFileSize,
	}, nil
}

func (writer *Writer) Close() error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if err := writer.closeFile(); err != nil {
		return fmt.Errorf("writer/file.Close: %v", err)
	}
	return nil
}

func (writer *Writer) Write(bs []byte, record *iface.Record) error {
	writer.lock.Lock()

	err := writer.checkFile(record)
	if err == nil {
		var n int
		n, err = writer.writer.Write(bs)
		writer.fileSize += int64(n)
	}

	writer.lock.Unlock()

	if err != nil {
		return fmt.Errorf("writer/file.Write: %v", err)
	}
	return nil
}

func (writer *Writer) Path() string {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.path
}

func (writer *Writer) SetPath(path string) error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if path == writer.path {
		return nil
	}
	if err := checkPath(path); err != nil {
		return fmt.Errorf("writer/file.SetPath: %v", err)
	}
	if err := writer.closeFile(); err != nil {
		return fmt.Errorf("writer/file.SetPath: %v", err)
	}
	writer.path = path
	return nil
}

func (writer *Writer) MaxFileSize() int64 {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.maxFileSize
}

func (writer *Writer) SetMaxFileSize(size int64) error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if size <= 0 {
		return errors.New("writer/file.SetMaxFileSize: size must be positive")
	}
	writer.maxFileSize = size
	return nil
}

func (writer *Writer) checkFile(record *iface.Record) error {
	if writer.writer == nil ||
		writer.day != record.Time.YearDay() ||
		writer.fileSize >= writer.maxFileSize {
		return writer.createFile(record)
	} else if time.Since(writer.checkTime) >= checkInterval {
		writer.checkTime = time.Now()
		_, err := os.Stat(writer.pathname)
		if os.IsNotExist(err) {
			return writer.createFile(record)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (writer *Writer) createFile(record *iface.Record) error {
	if err := writer.closeFile(); err != nil {
		return err
	}

	date := fmt.Sprintf(dateFormat, record.Time.Year(), record.Time.Month(), record.Time.Day())
	path := filepath.Join(writer.path, date)
	if err := os.MkdirAll(path, dirPerm); err != nil {
		return err
	}

	clock := fmt.Sprintf(timeFormat, record.Time.Hour(), record.Time.Minute(),
		record.Time.Second(), record.Time.Nanosecond()/1000)
	filename := clock + extension
	pathname := filepath.Join(path, filename)
	file, err := os.Create(pathname)
	if err != nil {
		return err
	}

	writer.writer = file
	writer.pathname = pathname
	writer.day = record.Time.YearDay()
	writer.fileSize = 0

	return nil
}

func (writer *Writer) closeFile() error {
	if writer.writer != nil {
		if err := writer.writer.Close(); err != nil {
			return err
		}
		writer.writer = nil
	}
	return nil
}

func checkPath(path string) error {
	if err := os.MkdirAll(path, dirPerm); err != nil {
		return err
	}
	return syscall.Access(path, 7 /* R_OK | W_OK | X_OK */)
}
