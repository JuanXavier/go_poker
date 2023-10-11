package p2p

import (
	// "io"
)

type Message struct {
	Payload any
	From    string
}

type Handshake struct {
	Version     string
	GameVariant GameVariant
	GameStatus  GameStatus
}

type MessagePeerList struct {
	Peers []string
}

func NewMessage(from string, payload any) *Message {
	return &Message{
		From:    from,
		Payload: payload,
	}
}
