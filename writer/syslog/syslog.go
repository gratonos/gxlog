package syslog

import (
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

type syslog struct {
	conn net.Conn
}

func syslogDial() (*syslog, error) {
	log := new(syslog)
	if err := log.connect(); err != nil {
		return nil, err
	}
	return log, nil
}

func (this *syslog) Write(timestamp time.Time, priority int, tag string, msg []byte) error {
	if this.conn != nil {
		if err := this.write(timestamp, priority, tag, msg); err == nil {
			return nil
		}
		this.Close()
	}
	if err := this.connect(); err != nil {
		return err
	}
	if err := this.write(timestamp, priority, tag, msg); err != nil {
		this.Close()
		return err
	}
	return nil
}

func (this *syslog) Close() error {
	if this.conn != nil {
		err := this.conn.Close()
		this.conn = nil
		return err
	}
	return nil
}

func (this *syslog) connect() error {
	networks := []string{"unixgram", "unix"}
	paths := []string{"/dev/log", "/var/run/syslog", "/var/run/log"}
	for _, network := range networks {
		for _, path := range paths {
			if conn, err := net.Dial(network, path); err == nil {
				this.conn = conn
				return nil
			}
		}
	}
	return errors.New("unix syslog delivery error")
}

func (this *syslog) write(timestamp time.Time, priority int, tag string, msg []byte) error {
	_, err := fmt.Fprintf(this.conn, "<%d>%s %s[%d]: %s",
		priority, timestamp.Format(time.Stamp), tag, os.Getpid(), msg)
	return err
}
