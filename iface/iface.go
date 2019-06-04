// Package iface contains basic definitions for logging.
package iface

import (
	"time"
)

type Level int

const (
	Trace Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
	Off
)

const LevelCount = Off - Trace

type Context struct {
	Key   string
	Value string
}

type Auxiliary struct {
	Prefix   string
	Contexts []Context
	Mark     bool
}

type Record struct {
	Time  time.Time
	Level Level
	File  string
	Line  int
	Pkg   string
	Func  string
	Msg   string
	Aux   Auxiliary
}

type Formatter interface {
	Format(record *Record) []byte
}

type Writer interface {
	Write(bs []byte, record *Record) error
}
