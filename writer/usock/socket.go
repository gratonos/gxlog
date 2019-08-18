package usock

import (
	"net"
	"sync"
)

type socket struct {
	listener net.Listener
	conns    map[int64]net.Conn
	id       int64
	wg       sync.WaitGroup
	lock     sync.Mutex
}

func openSocket(path string) (*socket, error) {
	listener, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	sock := &socket{
		listener: listener,
		conns:    make(map[int64]net.Conn),
	}
	sock.wg.Add(1)

	go sock.serve()

	return sock, nil
}

func (this *socket) Close() error {
	if err := this.listener.Close(); err != nil {
		return err
	}

	this.wg.Wait()

	this.lock.Lock()
	defer this.lock.Unlock()

	for id, conn := range this.conns {
		conn.Close()
		delete(this.conns, id)
	}

	return nil
}

func (this *socket) Write(bs []byte) {
	this.lock.Lock()

	for id, conn := range this.conns {
		if _, err := conn.Write(bs); err != nil {
			conn.Close()
			delete(this.conns, id)
		}
	}

	this.lock.Unlock()
}

func (this *socket) serve() {
	for {
		conn, err := this.listener.Accept()
		if err != nil {
			break
		}

		this.lock.Lock()

		id := this.id
		this.id++
		this.conns[id] = conn

		this.lock.Unlock()
	}

	this.wg.Done()
}
