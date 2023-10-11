package main

import (
	// "fmt"
	"log"
	// "github.com/juanxavier/go_poker/deck"
	"github.com/juanxavier/go_poker/p2p"
	"time"
)

func main() {
	/* ----------------------- SERVER ----------------------- */
	cfg := p2p.ServerConfig{
		Version:     "0.1.0",
		ListenAddr:  ":3000",
		GameVariant: p2p.TexasHoldEm,
	}
	server := p2p.NewServer(cfg)
	go server.Start()

	// Wait
	time.Sleep(1 * time.Second)

	/* -------------------- REMOTE SERVER ------------------- */
	remoteCfg := p2p.ServerConfig{
		Version:     "0.1.0",
		ListenAddr:  ":4000",
		GameVariant: p2p.TexasHoldEm,
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()

	if err := remoteServer.Connect(":3000"); err != nil {
		log.Fatal(err)
	}

	select {}

}
