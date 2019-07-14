package file

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
)

type Config struct {
	Dir         string
	MaxFileSize int64
}

func (this *Config) SetDefaults() {
	if this.Dir == "" {
		this.Dir = fmt.Sprintf("%s.%d", filepath.Base(os.Args[0]), os.Getpid())
	}

	if this.MaxFileSize == 0 {
		this.MaxFileSize = 20 * 1024 * 1024
	} else if this.MaxFileSize < 0 {
		this.MaxFileSize = math.MaxInt64
	}
}
