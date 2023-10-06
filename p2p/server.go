package p2p

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

/* ****************************************************** */
/*                          TYPES                         */
/* ****************************************************** */

type ServerConfig struct {
	Version    string
	ListenAddr string
}

type Server struct {
	ServerConfig
	handler   Handler
	transport *TCPTransport
	mu        sync.RWMutex // mutex for concurrent access,
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
}

/* ****************************************************** */
/*                        FUNCTIONS                       */
/* ****************************************************** */
// NewServer initializes and returns a new Server instance.
func NewServer(cfg ServerConfig) *Server {

	s := &Server{
		handler:      &DefaultHandler{},
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
	}
	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr
	tr.AddPeer = s.addPeer
	tr.DelPeer = s.addPeer
	return s
}

// starts the server and its associated routines.
func (s *Server) Start() {
	go s.loop()

	fmt.Printf("game server running on port %s\n", s.ListenAddr)

	logrus.WithFields(logrus.Fields{
		"port": s.ListenAddr,
		"type": "Texas Hold'em Poker",
	}).Info("Started new game server")

	s.transport.ListenAndAccept()
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}
	s.addPeer <- peer
	return peer.Send([]byte(s.Version))
}

func (s *Server) loop() {
	for {
		select {
		/* ------------------- DELETE PEER ------------------- */
		case peer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Player disconnected")

			delete(s.peers, peer.conn.RemoteAddr())

		/* ---------------------- ADD PEER ---------------------- */
		case peer := <-s.addPeer:
			go peer.ReadLoop(s.msgCh)

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("New player connected")

			s.peers[peer.conn.RemoteAddr()] = peer

			/* ----------------------- MESSAGE ---------------------- */
		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
