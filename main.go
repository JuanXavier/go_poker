package main

import (
	// "fmt"
	// "github.com/juanxavier/go_poker/deck"
	"github.com/juanxavier/go_poker/p2p"
	"time"
)

func main() {
	playerA := makeServerAndStart(":3000")
	playerB := makeServerAndStart(":4000")
	playerC := makeServerAndStart(":5000")
	// playerD := makeServerAndStart(":6000")
	// playerE := makeServerAndStart(":7000")
	// playerF := makeServerAndStart(":8000")

	playerB.Connect(playerA.ListenAddr)
	time.Sleep(200 * time.Millisecond)
	playerC.Connect(playerB.ListenAddr)
	time.Sleep(200 * time.Millisecond)
	// playerD.Connect(playerC.ListenAddr)
	// time.Sleep(200 * time.Millisecond)
	// playerE.Connect(playerD.ListenAddr)
	// time.Sleep(200 * time.Millisecond)
	// playerF.Connect(playerE.ListenAddr)

	select {}
}

func makeServerAndStart(addr string) *p2p.Server {
	cfg := p2p.ServerConfig{
		Version:     "0.1.0",
		ListenAddr:  addr,
		GameVariant: p2p.TexasHoldEm,
	}
	server := p2p.NewServer(cfg)
	go server.Start()
	time.Sleep(1 * time.Second)
	return server
}

func testPeersConnected(s *p2p.Server) {
}
