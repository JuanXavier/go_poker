package p2p

import (
	// "io"
	"net"
)

type Message struct {
	Payload any
	From    net.Addr
}

type Handshake struct {
	Version     string
	GameVariant GameVariant
	GameStatus  GameStatus
}

type MessagePeerList []net.Addr
