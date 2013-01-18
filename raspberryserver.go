package main
/*
Webserver implemented on RaspberryPi hardware
Collect weather data from a website/api (or another server masquerading as a sensor)
Moisture sensor gives soil moisture data
RaspberryPi decides whether or not to turn on watering system/valve

*/

import (
	"encoding/json"
	//"fmt"
	//"sync"
	//"time"
	"math/rand"
	"net/http"
)

type data struct {
	Moisture float64
	
}

func handle(w http.ResponseWriter, r *http.Request) {
	if err := json.NewEncoder(w).Encode(data{10*rand.Float64()}); err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", http.HandlerFunc(handle))
	if err := http.ListenAndServe(":25601", nil); err != nil {
		panic(err)
	}
}
