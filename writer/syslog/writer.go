// Package syslog implements a syslog writer.
package syslog

import (
	"fmt"
	"sync"

	"github.com/gratonos/gxlog/iface"
)

type Writer struct {
	tag      string
	facility Facility

	severities [iface.LevelCount]Severity
	log        *syslog
	lock       sync.Mutex
}

func Open(config Config) (*Writer, error) {
	config.SetDefaults()

	log, err := syslogDial()
	if err != nil {
		return nil, fmt.Errorf("writer/syslog.Open: %v", err)
	}

	writer := &Writer{
		tag:      config.Tag,
		facility: config.Facility,
		log:      log,
	}
	writer.MapSeverities(config.SeverityMap)
	return writer, nil
}

func (this *Writer) Close() error {
	this.lock.Lock()
	defer this.lock.Unlock()

	if err := this.log.Close(); err != nil {
		return fmt.Errorf("writer/syslog.Close: %v", err)
	}
	return nil
}

func (this *Writer) Write(bs []byte, record *iface.Record) error {
	this.lock.Lock()

	severity := this.severities[record.Level]
	priority := int(this.facility) | int(severity)
	err := this.log.Write(record.Time, priority, this.tag, bs)

	this.lock.Unlock()

	if err != nil {
		return fmt.Errorf("writer/syslog.Write: %v", err)
	}
	return nil
}

func (this *Writer) Tag() string {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.tag
}

func (this *Writer) SetTag(tag string) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.tag = tag
}

func (this *Writer) Facility() Facility {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.facility
}

func (this *Writer) SetFacility(facility Facility) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.facility = facility
}

func (this *Writer) Severity(level iface.Level) Severity {
	this.lock.Lock()
	defer this.lock.Unlock()

	return this.severities[level]
}

func (this *Writer) SetSeverity(level iface.Level, severity Severity) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.severities[level] = severity
}

func (this *Writer) MapSeverities(severityMap map[iface.Level]Severity) {
	this.lock.Lock()
	defer this.lock.Unlock()

	for level, severity := range severityMap {
		this.severities[level] = severity
	}
}
