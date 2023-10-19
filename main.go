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
	playerD := makeServerAndStart(":6000")

	time.Sleep(1 * time.Second)
	playerB.Connect(playerA.ListenAddr)
	time.Sleep(1 * time.Second)
	playerC.Connect(playerB.ListenAddr)
	playerD.Connect(playerC.ListenAddr)

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
