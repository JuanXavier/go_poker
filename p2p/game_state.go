package p2p

import (
	"fmt"
)

type GameStatus uint32

func (g GameStatus) String() string {
	switch g {
	case GameStatusWaiting:
		return "Waiting"
	case GameStatusDealing:
		return "Dealing"
	case GameStatusPreFlop:
		return "Pre-flop"
	case GameStatusFlop:
		return "Flop"
	case GameStatusTurn:
		return "Turn"
	case GameStatusRiver:
		return "River"
	default:
		return "Unknown"
	}
}

const (
	GameStatusWaiting GameStatus = iota
	GameStatusDealing
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

type GameState struct {
	isDealer   bool       // should be atomic accessible
	gameStatus GameStatus // should be atomic accessible
}

func NewGameState() *GameState {
	return &GameState{}
}

func (g *GameState) loop() {
	for {
		select {
		default:
			fmt.Println("Unknown message type")
		}
	}
}
