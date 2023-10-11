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

type MessagePeerList struct {
	Peers []net.Addr
}

func NewMessage(from net.Addr, payload any) *Message {
	return &Message{
		From:    from,
		Payload: payload,
	}
}
