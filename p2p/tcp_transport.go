package p2p

import (
	// "bytes"
	"encoding/gob"
	"github.com/sirupsen/logrus"
	"net"
)

type NetAddr string

func (n NetAddr) String() string {
	return string(n)
}

func (n NetAddr) Network() string {
	return "tcp"
}

type Peer struct {
	conn       net.Conn
	outbound   bool
	listenAddr string
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *Peer) ReadLoop(msgCh chan *Message) {
	for {
		msg := new(Message)
		if err := gob.NewDecoder(p.conn).Decode(&msg); err != nil {
			logrus.Errorf("decode msg error %s", err)
			break
		}

		msgCh <- msg
	}
	p.conn.Close()
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	AddPeer    chan *Peer
	DelPeer    chan *Peer
}

// func NewTCPTransport(addr string, addPeer chan *Peer, delPeer chan *Peer) *TCPTransport {
func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: addr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}

		peer := &Peer{
			conn:     conn,
			outbound: false,
		}

		t.AddPeer <- peer
	}
}
