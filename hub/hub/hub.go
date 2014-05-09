package main

import (
	"time"

	"./hub"
)

func main() {
	//1: instantiera hub
	h := hub.New()
	time.Sleep(time.Hour)
}
