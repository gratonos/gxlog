// Package unix implements a unix domain socket writer.
package unix

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gratonos/gxlog/iface"
)

const dirPerm = 0770

type Writer struct {
	socket *socket
}

func Open(path string) (*Writer, error) {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return nil, openError(err)
	}
	if err := checkAndRemove(path); err != nil {
		return nil, openError(err)
	}
	socket, err := openSocket(path)
	if err != nil {
		return nil, openError(err)
	}
	return &Writer{socket: socket}, nil
}

func (writer *Writer) Close() error {
	if err := writer.socket.Close(); err != nil {
		return fmt.Errorf("writer/unix.Close: %v", err)
	}
	return nil
}

func (writer *Writer) Write(bs []byte, _ *iface.Record) error {
	writer.socket.Write(bs)
	return nil
}

func checkAndRemove(path string) error {
	if _, err := os.Stat(path); err != nil {
		return nil
	}
	return os.Remove(path)
}

func openError(err error) error {
	return fmt.Errorf("writer/unix.Open: %v", err)
}
