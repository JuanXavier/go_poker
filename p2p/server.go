package p2p

import (
	"bytes"
	"fmt"
	"io"
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

type Message struct {
	Payload io.Reader
	From    net.Addr
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

type ServerConfig struct {
	ListenAddr string
}

type Server struct {
	ServerConfig
	handler  Handler
	listener net.Listener
	mu       sync.RWMutex // mutex for concurrent access,
	peers    map[net.Addr]*Peer
	addPeer  chan *Peer
	msgCh    chan *Message
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
		msgCh:        make(chan *Message),
	}
}

// starts the server and its associated routines.
func (s *Server) Start() {
	go s.loop()

	if err := s.listen(); err != nil {
		panic(err)
	}

	fmt.Printf("game server running on port %s\n", s.ListenAddr)

	s.acceptLoop()
}

// Method belonging to server struct
// The purpose of the acceptLoop() method is to continuously accept incoming connections,
// create Peer objects for each connection, send them to the s.addPeer channel,
// and handle each connection concurrently in separate goroutines using the handleConn() method.
func (s *Server) acceptLoop() {
	// Infinite loop to continuously accept incoming connections , similar to while (true) {}
	for {
		// accepts an incoming connection from the net.listener of the Server object s.
		// The Accept() method blocks until a connection is made,
		// and it returns the conn object representing the connection and an err if any error occurs.
		conn, err := s.listener.Accept()

		// if an error occurred while accepting the connection, the program will panic, terminating abruptly.
		if err != nil {
			panic(err)
		}

		// The conn object representing the accepted connection is assigned to the conn field of the Peer struct.
		// The & operator is used to take the address of the Peer struct, creating a pointer to it
		peer := &Peer{
			conn: conn,
		}

		// sends the peer object to the s.addPeer channel. The <- operator is used to send data into the channel.
		// The peer object is sent as a value,
		// which means it will be received as a value on the receiving end of the channel.
		s.addPeer <- peer

		peer.Send([]byte("go_poker_v0.1\n"))

		//starts a new goroutine by calling the handleConn() method of the Server object s and passing the conn
		// object as an argument. Goroutines allow concurrent execution of functions,
		//and this line enables the handleConn() method to run concurrently with the acceptLoop() method.
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}

		s.msgCh <- &Message{
			From:    conn.RemoteAddr(),
			Payload: bytes.NewReader(buf[:n]),
		}

		fmt.Println(string(buf[:n]))
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
			s.peers[peer.conn.RemoteAddr()] = peer
			fmt.Printf("new player connected %s\n", peer.conn.RemoteAddr())

		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
