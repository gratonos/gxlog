package logger

import (
	"github.com/gratonos/gxlog/iface"
)

type Config struct {
	Level  iface.Level
	Filter Filter
}

func (this *Config) SetDefaults() {
	this.Filter = fillFilter(this.Filter)
}
