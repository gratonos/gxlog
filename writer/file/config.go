package file

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
)

type Config struct {
	Path        string
	MaxFileSize int64
}

func (config *Config) SetDefaults() {
	if config.Path == "" {
		config.Path = fmt.Sprintf("%s.%d", filepath.Base(os.Args[0]), os.Getpid())
	}

	if config.MaxFileSize == 0 {
		config.MaxFileSize = 20 * 1024 * 1024
	} else if config.MaxFileSize < 0 {
		config.MaxFileSize = math.MaxInt64
	}
}
