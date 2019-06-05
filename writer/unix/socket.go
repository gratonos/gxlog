package unix

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

func (sock *socket) Close() error {
	if err := sock.listener.Close(); err != nil {
		return err
	}

	sock.wg.Wait()

	sock.lock.Lock()
	defer sock.lock.Unlock()

	for id, conn := range sock.conns {
		conn.Close()
		delete(sock.conns, id)
	}

	return nil
}

func (sock *socket) Write(bs []byte) {
	sock.lock.Lock()

	for id, conn := range sock.conns {
		if _, err := conn.Write(bs); err != nil {
			conn.Close()
			delete(sock.conns, id)
		}
	}

	sock.lock.Unlock()
}

func (sock *socket) serve() {
	for {
		conn, err := sock.listener.Accept()
		if err != nil {
			break
		}

		sock.lock.Lock()

		id := sock.id
		sock.id++
		sock.conns[id] = conn

		sock.lock.Unlock()
	}
	sock.wg.Done()
}
