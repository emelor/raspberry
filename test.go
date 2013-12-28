package main
import (
  "pi"
  "hub"
)

func main () {
//1: instantiera hub
hub := new(hub.Hub)
//hub := &hub.Hub{}
//2: instantiera pi
pi := new(pi.Pi)

//3: pi: connect to hub
pi.ConnectTo(hub)
hub.Register(pi)
pi.GetWeather()
hub.PushSettings(pi)
pi.Routine()
}

