package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

type Server interface {
	ListenOnPort(int) error
	ListenAnywhere() error
	Port() int
	Stop() error
}

func NewServer(readDeadline time.Duration, bytesCh chan []byte) Server {
	return &udpServer{
		bytesCh:      bytesCh,
		readDeadline: readDeadline,
	}
}

type udpServer struct {
	bytesCh chan []byte

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
	addr, err := net.ResolveUDPAddr("udp", "localhost:0")
	if err != nil {
		return err
	}

	return u.start(addr)
}

func (u *udpServer) start(rawAddr *net.UDPAddr) error {
	conn, err := net.ListenUDP("udp", rawAddr)
	if err != nil {
		return err
	}

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return errors.New("invalid connection")
	}

	go func(conn *net.UDPConn) {
		// create a buffer which can store the output of the
		// largest acceptable, configured udp datagram
		buf := make([]byte, MaxUDPPayloadSize)

		for {
			conn.SetReadDeadline(time.Now().Add(u.readDeadline))

			bytesRead, _, err := conn.ReadFromUDP(buf)
			if err == nil {
				u.bytesCh <- buf[:bytesRead]
				continue
			}

			if timeoutErr, ok := err.(net.Error); ok && timeoutErr.Timeout() {
				continue
			}

			log.Println(err)
		}
	}(conn)

	u.conn = conn
	u.port = addr.Port
	return nil
}

func (u *udpServer) Port() int {
	if u.conn == nil {
		log.Panicf("server not started")
	}

	return u.port
}

func (u *udpServer) Stop() error {
	if err := u.conn.Close(); err != nil {
		return err
	}
	return nil
}
