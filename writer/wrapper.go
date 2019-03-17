package writer

import (
	"io"

	"github.com/gratonos/gxlog/iface"
)

func Wrap(writer io.Writer) iface.Writer {
	return Func(func(bs []byte, _ *iface.Record) error {
		_, err := writer.Write(bs)
		return err
	})
}
