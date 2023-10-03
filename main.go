package main

import (
	"fmt"
	"github.com/juanxavier/go_poker/deck"
) 

func main () {
	// card := deck.NewCard(deck.Spades, 3)
	d := deck.New()
	fmt.Println(d)
}