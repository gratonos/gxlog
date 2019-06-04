package syslog

import (
	"os"
	"path/filepath"

	"github.com/gratonos/gxlog/iface"
)

type Facility int

const (
	FacKern Facility = iota << 3
	FacUser
	FacMail
	FacDaemon
	FacAuth
	FacSyslog
	FacLPR
	FacNews
	FacUUCP
	FacCron
	FacAuthPriv
	FacFTP
)

type Severity int

const (
	SevEmerg Severity = iota
	SevAlert
	SevCrit
	SevErr
	SevWarning
	SevNotice
	SevInfo
	SevDebug
)

type Config struct {
	Tag         string
	Facility    Facility
	SeverityMap map[iface.Level]Severity
}

func (config *Config) SetDefaults() {
	if config.Tag == "" {
		config.Tag = filepath.Base(os.Args[0])
	}

	severityMap := map[iface.Level]Severity{
		iface.Trace: SevDebug,
		iface.Debug: SevDebug,
		iface.Info:  SevInfo,
		iface.Warn:  SevWarning,
		iface.Error: SevErr,
		iface.Fatal: SevCrit,
	}
	for level, severity := range config.SeverityMap {
		severityMap[level] = severity
	}
	config.SeverityMap = severityMap
}
