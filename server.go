package main

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Server interface {
	ListenOnPort(int) error
	ListenAnywhere() error
	Port() int
	Stop() error
}

func NewServer(readDeadline time.Duration, receiverCh chan []byte) {
	return &udpServer{
		receiverCh:   receiverCh,
		readDeadline: readDeadline,
	}
}

type udpServer struct {
	receiverCh chan []byte

	// udp listener config
	port         int
	conn         *net.UDPConn
	readDeadline time.Duration
}

func (u *udpServer) ListenOnPort(port int) error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	return u.start(addr)
}

func (u *udpServer) ListenAnywhere() error {
	addr, err := net.ResolveUDPAddr("udp", "localhost")
	if err != nil {
		return err
	}

	return nil
}

func (u *udpServer) start(rawAddr *net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", rawAddr)
	if err != nil {
		return err
	}

	addr, ok := conn.LAddr().(*net.UDPAddr)
	if !ok {
		return errors.New("invalid connection")
	}

	go func(conn *net.UDPConn) {
		// create a buffer which can store the output of the
		// largest acceptable, configured udp datagram
		buf := make([]byte, 0, MAXUDPPayloadSize)

		for {
			conn.SetReadDeadline(time.Now().Add(u.readDeadline))

			bytesRead, originAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				// TODO: parse out errors that are related to
				// the connection being closed and from a bad
				// actor
			}

			// send the udp packets to a goroutine to be processed
			// and parsed as a message
			bytesCh <- buf[:bytesRead]
		}
	}()

	u.port = addr.Port
	return nil
}

func (u *udpServer) Stop() error {
	if err := u.conn.Close(); err != nil {
		return err
	}
	return nil
}
