package main

import (
	"fmt"
	"github.com/juanxavier/go_poker/deck"
)

func main() {
	for i := 0; i < 10; i++ {
		d := deck.New()
		fmt.Println(d)
	}
}
