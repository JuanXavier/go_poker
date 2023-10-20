package p2p

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type GameStatus uint32

func (g GameStatus) String() string {
	switch g {
	case GameStatusWaitingForCards:
		return "WAITING FOR CARDS"
	case GameStatusDealing:
		return "DEALING"
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
	GameStatusWaitingForCards GameStatus = iota
	GameStatusDealing
	GameStatusPreFlop
	GameStatusFlop
	GameStatusTurn
	GameStatusRiver
)

type Player struct {
	Status GameStatus
}

type GameState struct {
	isDealer   bool       // should be atomic accessible
	gameStatus GameStatus // should be atomic accessible

	playersLock sync.RWMutex
	players     map[string]*Player
}

func (g *GameState) AddPlayer(addr string, status GameStatus) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()

	g.players[addr] = &Player{
		Status: status,
	}

	logrus.WithFields(logrus.Fields{
		"addr":   addr,
		"status": status,
	}).Info("New player joined")
}

func NewGameState() *GameState {
	g := &GameState{
		isDealer:   false,
		gameStatus: GameStatusWaitingForCards,
		players:    make(map[string]*Player),
	}
	go g.loop()

	return g
}

func (g *GameState) loop() {

	ticker := time.NewTicker(time.Second * 5)

	for {

		select {
		case <-ticker.C:

			logrus.WithFields(logrus.Fields{
				"connected players": len(g.players),
				"status":            g.gameStatus,
			}).Info("New player joined")

		default:
			fmt.Println("Unknown message type")
		}
	}
}
