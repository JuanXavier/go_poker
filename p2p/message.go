package p2p

import (
	"github.com/juanxavier/go_poker/deck"
)

type Message struct {
	Payload any
	From    string
}

type Handshake struct {
	Version     string
	GameVariant GameVariant
	GameStatus  GameStatus
	ListenAddr  string
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

type MessageCards struct {
	Deck deck.Deck
}
