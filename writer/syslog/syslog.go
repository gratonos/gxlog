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

func (log *syslog) Write(timestamp time.Time, priority int, tag string, msg []byte) error {
	if log.conn != nil {
		if err := log.write(timestamp, priority, tag, msg); err == nil {
			return nil
		} else {
			log.Close()
		}
	}
	if err := log.connect(); err != nil {
		return err
	}
	if err := log.write(timestamp, priority, tag, msg); err != nil {
		log.Close()
		return err
	}
	return nil
}

func (log *syslog) Close() error {
	if log.conn != nil {
		err := log.conn.Close()
		log.conn = nil
		return err
	}
	return nil
}

func (log *syslog) connect() error {
	networks := []string{"unixgram", "unix"}
	paths := []string{"/dev/log", "/var/run/syslog", "/var/run/log"}
	for _, network := range networks {
		for _, path := range paths {
			if conn, err := net.Dial(network, path); err == nil {
				log.conn = conn
				return nil
			}
		}
	}
	return errors.New("Unix syslog delivery error")
}

func (log *syslog) write(timestamp time.Time, priority int, tag string, msg []byte) error {
	_, err := fmt.Fprintf(log.conn, "<%d>%s %s[%d]: %s",
		priority, timestamp.Format(time.Stamp), tag, os.Getpid(), msg)
	return err
}
