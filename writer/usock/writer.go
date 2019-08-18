// Package usock implements a unix domain socket writer.
package usock

import (
	"fmt"
	"os"
	"path/filepath"

	gos "github.com/gratonos/goutil/os"
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
	if err := gos.RemoveIfExists(path); err != nil {
		return nil, openError(err)
	}
	socket, err := openSocket(path)
	if err != nil {
		return nil, openError(err)
	}
	return &Writer{socket: socket}, nil
}

func (this *Writer) Close() error {
	if err := this.socket.Close(); err != nil {
		return fmt.Errorf("writer/usock.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, _ *iface.Record) error {
	this.socket.Write(bs)
	return nil
}

func openError(err error) error {
	return fmt.Errorf("writer/usock.Open: %v", err)
}
