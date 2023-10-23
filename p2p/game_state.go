package p2p

import (
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

type GameStatus int32

const (
	GameStatusWaitingForCards GameStatus = iota
	GameStatusReceivingCards
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
	isDealer               bool       // should be atomic accessible
	gameStatus             GameStatus // should be atomic accessible
	broadcast              chan BroadcastTo
	playersWaitingForCards int32
	playersLock            sync.RWMutex
	players                map[string]*Player
	listenAddr             string
}

/* ---------------------- FUNCTIONS --------------------- */
func (g GameStatus) String() string {
	switch g {
	case GameStatusWaitingForCards:
		return "WAITING FOR CARDS"
	case GameStatusReceivingCards:
		return "RECEIVING CARDS"
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

/* -------------------------- - ------------------------- */

func (g *GameState) SendToPlayerWithStatus(payload any, s GameStatus) {
	players := g.GetPlayersWithStatus(s)

	g.broadcast <- BroadcastTo{
		To:      players,
		Payload: payload,
	}
}

/* -------------------------- - ------------------------- */

func (g *GameState) GetPlayersWithStatus(s GameStatus) []string {
	players := []string{}
	for addr, _ := range g.players {
		players = append(players, addr)
	}
	return players
}

/* -------------------------- - ------------------------- */

func (g *GameState) AddPlayerWaitingForCards() {
	atomic.AddInt32(&g.playersWaitingForCards, 1)
}

/* -------------------------- - ------------------------- */

func (g *GameState) CheckNeedDealCards() {
	playersWaiting := atomic.LoadInt32(&g.playersWaitingForCards)

	if playersWaiting == int32(g.LenPlayersConnectedWithLock()) && g.isDealer && g.gameStatus == GameStatusWaitingForCards {

		// do something
		logrus.WithFields(logrus.Fields{
			"addr": g.listenAddr,
		}).Info("Need to deal cards")

		g.InitiateShuffleAndDeal()
	}
}

/* -------------------------- - ------------------------- */

func (g *GameState) DealCards() {
	// g.broadcast <- MessageEncDeck{}
}

/* -------------------------- - ------------------------- */

func (g *GameState) SetPlayerStatus(addr string, status GameStatus) {
	player, ok := g.players[addr]

	if !ok {
		panic("Player not found, although it should exist")
	}

	player.Status = status
	g.CheckNeedDealCards()
}

/* -------------------------- - ------------------------- */

func (g *GameState) LenPlayersConnectedWithLock() int {
	g.playersLock.RLock()
	defer g.playersLock.RUnlock()
	return len(g.players)
}

/* -------------------------- - ------------------------- */

func (g *GameState) AddPlayer(addr string, status GameStatus) {
	g.playersLock.Lock()
	defer g.playersLock.Unlock()

	if status == GameStatusWaitingForCards {
		g.AddPlayerWaitingForCards()
	}

	g.players[addr] = new(Player)

	// set the player status when adding the player
	g.SetPlayerStatus(addr, status)

	logrus.WithFields(logrus.Fields{
		"addr":   addr,
		"status": status,
	}).Info("New player joined")
}

/* -------------------------- - ------------------------- */

func NewGameState(addr string, broadcast chan BroadcastTo) *GameState {
	g := &GameState{
		listenAddr: addr,
		broadcast:  broadcast,
		isDealer:   false,
		gameStatus: GameStatusWaitingForCards,
		players:    make(map[string]*Player),
	}
	go g.loop()
	return g
}

// todo check other RW occurrences of the GameStatus
func (g *GameState) SetStatus(s GameStatus) {
	// Only update when status is different
	if g.gameStatus != s {
		atomic.StoreInt32((*int32)(&g.gameStatus), (int32)(s))
	}
}

func (g *GameState) loop() {
	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			logrus.WithFields(logrus.Fields{
				"we":                g.listenAddr,
				"connected players": g.LenPlayersConnectedWithLock(),
				"status":            g.gameStatus,
			}).Info("New player joined")

		default:
			// logrus.Info("Unknown error type")
		}
	}
}

// only used for the "real" dealer
func (g *GameState) InitiateShuffleAndDeal() {
	g.SetStatus(GameStatusReceivingCards)
	// g.broadcast <- MessageEncDeck{Deck: [][]byte{}}
	g.SendToPlayerWithStatus(MessageEncDeck{Deck: [][]byte{}}, GameStatusWaitingForCards)
}
