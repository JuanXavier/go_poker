package p2p

import (
	"github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
	"time"
)

type Player struct {
	Status     GameStatus
	ListenAddr string
}

type GameState struct {
	listenAddr             string
	broadcast              chan BroadcastTo
	isDealer               bool       // should be atomic accessible
	gameStatus             GameStatus // should be atomic accessible
	playersWaitingForCards int32
	players                map[string]*Player
	playersLock            sync.RWMutex
	playersList            map[string]*Player
	decksReceived          map[string]bool
	decksReceivedLock      sync.RWMutex
}

/* ---------------------- FUNCTIONS --------------------- */

/* -------------------------- - ------------------------- */

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

	player := &Player{ListenAddr: addr}
	g.players[addr] = player
	g.playersList = append(g.playersList, player)

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
		listenAddr:    addr,
		broadcast:     broadcast,
		isDealer:      false,
		gameStatus:    GameStatusWaitingForCards,
		players:       make(map[string]*Player),
		decksReceived: make(map[string]bool),
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
				"decksReceived":     g.decksReceived,
			}).Info("New player joined")

		default:
			// logrus.Info("Unknown error type")
		}
	}
}

func (g *GameState) SetDecksReceived(from string) {
	g.decksReceivedLock.Lock()
	g.decksReceived[from] = true
	g.decksReceivedLock.Unlock()
}

func (g *GameState) ShuffleAndEncrypt(from string, deck [][]byte) error {
	//TODO
	dealToPlayer := g.playersList[1]
	return nil
}

// only used for the "real" dealer
func (g *GameState) InitiateShuffleAndDeal() {
	dealToPlayer := g.playersList["0"]
	g.SendToPlayer(dealToPlayer.ListenAddr, MessageEncDeck{Deck: [][]byte{}})
	g.SetStatus(GameStatusReceivingCards)
}

func (g *GameState) SendToPlayer(addr string, payload any) error {
	g.broadcast <- BroadcastTo{
		To:      []string{addr},
		Payload: payload,
	}
	logrus.WithFields(logrus.Fields{
		"payload": payload,
		"player":  addr,
	}).Info("Sending payload to players")
	return nil
}

func (g *GameState) SendToPlayerWithStatus(payload any, s GameStatus) {
	// players := g.GetPlayersWithStatus(s)

	// g.broadcast <- BroadcastTo{
	// 	To:      players,
	// 	Payload: payload,
	// }
	// logrus.WithFields(logrus.Fields{
	// 	"payload": payload,
	// 	"players": players,
	// }).Info("Sending payload to players")
}
