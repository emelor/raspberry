package main

import (
	"fmt"
	"time"

	"./hub"
	"./pi"
)

func main() {
	//1: instantiera hub
	h := hub.New()
	//2: instantiera pi
	p := pi.New()
	fmt.Println(p, h)

	//3: pi: connect to hub
	p.ConnectTo(h)
	time.Sleep(time.Hour)
}
