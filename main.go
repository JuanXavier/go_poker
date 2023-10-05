package main

import (
	"fmt"
	// "github.com/juanxavier/go_poker/deck"
	"github.com/juanxavier/go_poker/p2p"
	"time"
)

func main() {
	/* ----------------------- SERVER ----------------------- */
	cfg := p2p.ServerConfig{
		Version:    "0.1.0",
		ListenAddr: ":3001",
	}
	server := p2p.NewServer(cfg)
	go server.Start()
	time.Sleep(1 * time.Second)

	/* -------------------- REMOTE SERVER ------------------- */
	remoteCfg := p2p.ServerConfig{
		Version:    "0.1.0",
		ListenAddr: ":4000",
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()

	if err := remoteServer.Connect(":3001"); err != nil {
		fmt.Println(err)
	}

	select {}

}
