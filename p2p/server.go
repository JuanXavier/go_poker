package server

import (
	"fmt"
	"net"
	"sync"
)

/* ****************************************************** */
/*                          TYPES                         */
/* ****************************************************** */

type Peer struct {
	conn net.Conn
}

type TCPTransport struct {
}

type ServerConfig struct {
	ListenAddr string
}

type Server struct {
	ServerConfig
	listener net.Listener
	mu       sync.RWMutex
	peers    map[net.Addr]*Peer
	addPeer  chan *Peer
}

/* ****************************************************** */
/*                        FUNCTIONS                       */
/* ****************************************************** */
// NewServer initializes and returns a new Server instance.
func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
	}
}

// starts the server and its associated routines.
func (s *Server) Start() {
	go s.loop()

	if err := s.listen(); err != nil {
		panic(err)
	}
	conn, err := s.listener.Accept()
	if err != nil {
		panic(err)
	}
	go s.handleConn(conn)
}

func (s *Server) handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Println(string(buf[:n]))
	}
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConn(conn)
	}
}

func (s *Server) listen() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.listener = ln
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			fmt.Printf("new player connected %s", peer.conn.RemoteAddr())
		}
	}
}
