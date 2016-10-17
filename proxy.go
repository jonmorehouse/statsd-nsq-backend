package main

import (
	"log"
	"net"
	"sync"
	"time"
)

type Proxier interface {
	Proxy([]byte) error
	Stop() error
}

// MultiProxier is a UDP
func NewMultiProxier(addrs []string) Proxier {
	proxiers := make([]Proxier, len(addrs))
	for idx, addr := range addrs {
		proxiers[idx] = NewProxier(addr)
	}

	return &multiProxier{
		proxiers: proxiers,
	}
}

func (m *multiProxier) Stop() error {
	var wg sync.WaitGroup

	for _, proxier := range m.proxiers {
		wg.Add(1)

		go func(proxier Proxier) {
			if err := proxier.Stop(); err != nil {
				log.Println(err)
			}
			wg.Done()
		}(proxier)
	}

	wg.Wait()
	return nil
}

func (m *multiProxier) Proxy(msg []byte) error {
	var wg sync.WaitGroup

	for _, proxier := range m.proxiers {
		wg.Add(1)

		go func(proxier Proxier) {
			if err := proxier.Proxy(msg); err != nil {
				log.Println(err)
			}
			wg.Done()
		}(proxier)
	}

	wg.Wait()
	return nil
}

type multiProxier struct {
	proxiers []Proxier
}

// NewProxier is a UDP proxy which writes messages via UDP
func NewProxier(addr string, writeTimeout time.Duration) (Proxier, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	return &proxier{
		addr: udpAddr,
		conn: conn,
	}, nil
}

type proxier struct {
	addr *net.UDPAddr
	conn *net.UDPConn
}

func (p *proxier) Stop() error {
	return p.conn.Close()
}

func (p *proxier) Proxy(msg []byte) error {
	//
}
