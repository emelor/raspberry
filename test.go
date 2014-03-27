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
	//Wait one second for hub to start server in separate go routine
	time.Sleep(time.Second)
	remoteHub, err := pi.NewRemoteHub(h.GetWSAddress())
	if err != nil {
		panic(err)
	}
	p.ConnectTo(remoteHub)
	time.Sleep(time.Hour)
}
