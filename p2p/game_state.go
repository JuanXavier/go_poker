package p2p

import (
	"fmt"
)

type Round uint32

const (
	Dealing Round = iota
	PreFlop
	Flop
	Turn
	River
)

type GameState struct {
	isDealer bool   // atomic accesible
	Round    uint32 // atomic accesible
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
