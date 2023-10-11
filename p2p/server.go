package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
)

/* ****************************************************** */
/*                          TYPES                         */
/* ******************************************************	 */

type GameVariant uint8

func (gv GameVariant) String() string {

	switch gv {
	case TexasHoldEm:
		return "Texas Hold'em"
	default:
		return "Unknown"
	}
}

const (
	TexasHoldEm GameVariant = iota
	Other
)

type ServerConfig struct {
	Version     string
	ListenAddr  string
	GameVariant GameVariant
}

type Server struct {
	ServerConfig
	transport *TCPTransport
	mu        sync.RWMutex // mutex for concurrent access,
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
	gameState *GameState
}

/* ****************************************************** */
/*                        FUNCTIONS                       */
/* ****************************************************** */
// NewServer initializes and returns a new Server instance.
func NewServer(cfg ServerConfig) *Server {

	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
		gameState:    NewGameState(),
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
		"port":    s.ListenAddr,
		"type":    "Texas Hold'em Poker",
		"variant": s.GameVariant,
	}).Info("Started new game server")

	s.transport.ListenAndAccept()
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn:     conn,
		outbound: true,
	}
	s.addPeer <- peer

	return s.SendHandshake(peer)
}

func (s *Server) loop() {
	for {
		select {
		/* ------------------- DELETE PEER ------------------- */
		case peer := <-s.delPeer:
			// Log
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Player disconnected")

			// Remove peer from map
			delete(s.peers, peer.conn.RemoteAddr())

			/* ---------------------- ADD PEER ---------------------- */
			// If a new peer connects to the server we send our handshake
			//message and wait for the response..
		case peer := <-s.addPeer:
			// Check for errors and log them
			if err := s.handshake(peer); err != nil {
				logrus.Errorf("Handshake with peer failed: %s", err)
				peer.conn.Close()
				continue
			}

			// TODO
			go peer.ReadLoop(s.msgCh)

			if !peer.outbound {
				if err := s.SendHandshake(peer); err != nil {
					logrus.Errorf("Failed to send handshake with peer: %s", err)
				}
			}

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Handshake successful: new player connected")

			s.peers[peer.conn.RemoteAddr()] = peer

		/* ----------------------- MESSAGE ---------------------- */
		case msg := <-s.msgCh:
			if err := s.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

type Handshake struct {
	Version     string
	GameVariant GameVariant
}

func (s *Server) SendHandshake(p *Peer) error {
	hs := &Handshake{
		GameVariant: s.GameVariant,
		Version:     s.Version,
	}

	buf := new(bytes.Buffer)

	// if err := hs.Encode(buf); err != nil {
	// 	return err
	// }

	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}

	return p.Send(buf.Bytes())
}

/* ****************************************************** */
/*                   ENCODING / DECODING                  */
/* ****************************************************** */

/* ****************************************************** */
/*                        HANDSHAKE                       */
/* ****************************************************** */
func (s *Server) handshake(p *Peer) error {
	hs := &Handshake{}

	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	// if err := hs.Decode(p.conn); err != nil {
	// 	return err
	// }

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("Invalid game variant %s", hs.GameVariant)
	}

	if s.Version != hs.Version {
		return fmt.Errorf("Invalid game version %s", hs.Version)

	}

	logrus.WithFields(logrus.Fields{
		"peer":    p.conn.RemoteAddr(),
		"version": hs.Version,
		"variant": hs.GameVariant,
	}).Info("New player connected")

	return nil
}

func (s *Server) HandleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	// return s.handler.HandleMessage(msg)
	return nil
}
