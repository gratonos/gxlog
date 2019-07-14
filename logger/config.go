package logger

import (
	"github.com/gratonos/gxlog/iface"
)

type Config struct {
	Level  iface.Level
	Filter Filter
}

func (this *Config) SetDefaults() {
	if this.Filter == nil {
		this.Filter = nullFilter
	}
}
