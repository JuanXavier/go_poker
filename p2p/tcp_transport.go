package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type Message struct {
	Payload io.Reader
	From    net.Addr
}

type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *Peer) ReadLoop(msgCh chan *Message) {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)

		if err != nil {
			break
		}

		msgCh <- &Message{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}

		fmt.Println(string(buf[:n]))
	}

	p.conn.Close()
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
}

func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: addr,
	}
}
