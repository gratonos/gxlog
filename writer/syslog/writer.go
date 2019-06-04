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

func (writer *Writer) Close() error {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	if err := writer.log.Close(); err != nil {
		return fmt.Errorf("writer/syslog.Close: %v", err)
	}
	return nil
}

func (writer *Writer) Write(bs []byte, record *iface.Record) error {
	writer.lock.Lock()

	severity := writer.severities[record.Level]
	priority := int(writer.facility) | int(severity)
	err := writer.log.Write(record.Time, priority, writer.tag, bs)

	writer.lock.Unlock()

	if err != nil {
		return fmt.Errorf("writer/syslog.Write: %v", err)
	}
	return nil
}

func (writer *Writer) Tag() string {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.tag
}

func (writer *Writer) SetTag(tag string) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.tag = tag
}

func (writer *Writer) Facility() Facility {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.facility
}

func (writer *Writer) SetFacility(facility Facility) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.facility = facility
}

func (writer *Writer) Severity(level iface.Level) Severity {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	return writer.severities[level]
}

func (writer *Writer) SetSeverity(level iface.Level, severity Severity) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	writer.severities[level] = severity
}

func (writer *Writer) MapSeverities(severityMap map[iface.Level]Severity) {
	writer.lock.Lock()
	defer writer.lock.Unlock()

	for level, severity := range severityMap {
		writer.severities[level] = severity
	}
}
