package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"time"
)

/* ****************************************************** */
/*                          TYPES                         */
/* ******************************************************	 */

type GameVariant uint8

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

const (
	TexasHoldEm GameVariant = iota
	Other
)

/* ****************************************************** */
/*                      GAME VARIANT                      */
/* ****************************************************** */
func (gv GameVariant) String() string {
	switch gv {
	case TexasHoldEm:
		return "Texas Hold'em"
	default:
		return "Unknown"
	}
}

/* ****************************************************** */
/*                         SERVER                         */
/* ****************************************************** */

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer, 10),
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

func (s *Server) Start() {
	go s.loop()

	logrus.WithFields(logrus.Fields{
		"port":    s.ListenAddr,
		"variant": s.GameVariant,
	}).Info("Started new game server")

	s.transport.ListenAndAccept()
}

func (s *Server) isInPeerList(addr string) bool {
	for _, peer := range s.peers {
		if peer.listenAddr == addr {
			return true
		}
	}
	return false
}

func (s *Server) Connect(addr string) error {
	if s.isInPeerList(addr) {
		return nil
	}

	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)

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
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("Player disconnected")

			delete(s.peers, peer.conn.RemoteAddr())

		/* ---------------------- ADD PEER ---------------------- */
		// If a new peer connects to the server we send our handshake
		//message and wait for the response..
		case peer := <-s.addPeer:
			if err := s.handleNewPeer(peer); err != nil {
				logrus.Errorf("handle peer error: %s", err)
			}
		/* ----------------------- MESSAGE ---------------------- */
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

/* ****************************************************** */
/*                        HANDSHAKE                       */
/* ****************************************************** */

func (s *Server) SendHandshake(p *Peer) error {
	hs := &Handshake{
		GameVariant: s.GameVariant,
		Version:     s.Version,
		GameStatus:  s.gameState.gameStatus,
		ListenAddr:  s.ListenAddr,
	}
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}
	return p.Send(buf.Bytes())
}

func (s *Server) handshake(p *Peer) (*Handshake, error) {
	hs := &Handshake{}

	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return nil, err
	}
	if s.GameVariant != hs.GameVariant {
		return nil, fmt.Errorf("Game variant does not match %s", hs.GameVariant)
	}
	if s.Version != hs.Version {
		return nil, fmt.Errorf("Invalid game version %s", hs.Version)
	}

	p.listenAddr = hs.ListenAddr

	return hs, nil
}

func (s *Server) handleMessage(msg *Message) error {

	switch v := msg.Payload.(type) {
	case MessagePeerList:
		return s.handlePeerList(v)
	}
	return nil
}

func init() {
	gob.Register(MessagePeerList{})
}

/* ****************************************************** */
/*                        PEER LIST                       */
/* ****************************************************** */

func (s *Server) handleNewPeer(peer *Peer) error {
	hs, err := s.handshake(peer)

	if err != nil {
		peer.conn.Close()
		delete(s.peers, peer.conn.RemoteAddr())
		return fmt.Errorf("%s: Handshake with peer failed: %s ", s.ListenAddr, err)
	}

	// always needs to start after the handshake
	go peer.ReadLoop(s.msgCh)

	if !peer.outbound {
		if err := s.SendHandshake(peer); err != nil {
			peer.conn.Close()
			delete(s.peers, peer.conn.RemoteAddr())
			return fmt.Errorf("Failed to send handshake with peer: %s", err)
		}

		if err := s.sendPeerList(peer); err != nil {
			return fmt.Errorf("Peer list error: %s", err)
		}
	}

	logrus.WithFields(logrus.Fields{
		"peer":       peer.conn.RemoteAddr(),
		"version":    hs.Version,
		"variant":    hs.GameVariant,
		"gameStatus": hs.GameStatus,
		"listenAddr": peer.listenAddr,
		"we":         s.ListenAddr,
	}).Info("Handshake successful: new player connected")

	//add to map
	s.peers[peer.conn.RemoteAddr()] = peer
	return nil
}

func (s *Server) handlePeerList(l MessagePeerList) error {

	logrus.WithFields(logrus.Fields{
		"we":   s.ListenAddr,
		"list": l.Peers,
	}).Info("received peer list message")

	fmt.Printf("peerList => %+v\n", l)

	for i := 0; i < len(l.Peers); i++ {
		if err := s.Connect(l.Peers[i]); err != nil {
			logrus.Errorf("Failed to connect to peer: %s", err)
			continue
		}
	}

	return nil
}

func (s *Server) sendPeerList(p *Peer) error {
	peerList := MessagePeerList{
		Peers: []string{},
	}

	for _, peer := range s.peers {
		peerList.Peers = append(peerList.Peers, peer.listenAddr)
	}

	if len(peerList.Peers) == 0 {
		return nil
	}

	msg := NewMessage(s.ListenAddr, peerList)

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}
	return p.Send(buf.Bytes())
}
