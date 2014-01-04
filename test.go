package main

import (
	"./hub"
	"./pi"
)

func main() {
	//1: instantiera hub
	h := hub.New()
	//2: instantiera pi
	p := pi.New()

	//3: pi: connect to hub
	p.ConnectTo(h)
}
