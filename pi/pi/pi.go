package main

import (
	"flag"
	"fmt"
	"time"
	"./pi"
)

func main() {
	server := flag.String("server", "", "Where to connect")
	flag.Parse()
	if *server == "" {
		flag.Usage()
		return
	}

	//2: instantiera pi
	p := pi.New()
	fmt.Println(p, h)

	//3: pi: connect to hub
	//Wait one second for hub to start server in separate go routine
	time.Sleep(time.Second)
	remoteHub, err := pi.NewRemoteHub(*server)
	if err != nil {
		panic(err)
	}
	p.ConnectTo(remoteHub)
	time.Sleep(time.Hour)
}
