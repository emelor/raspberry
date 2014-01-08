package pi

import (
	"fmt"
	"math/rand"
	"time"

	"../common"
)

type Pi struct {
	//
	config  common.Configuration
	weather common.Weather
	rain    float64
}

func New() *Pi {
	return &Pi{}

}

func (self *Pi) getMoisture() float64 {

	moisture := rand.Float64()
	return moisture
}

func (self *Pi) UpdateConfig(config common.Configuration) {
	self.config = config
	//save new config to file
}

func (self *Pi) UpdateWeather(weather common.Weather) {
	self.weather = weather
	//save new config to file
}

func (self *Pi) runPump() {
	fmt.Println("Pump running")

}

func (self *Pi) RoutineCheck() {
	for {
		var watering bool
		//Soil moisture sensor reading
		moisture := self.getMoisture()

		if (moisture < self.config.MoistureThreshold) && (self.weather.Rain <= self.config.RainLimit) {
			watering = true
			fmt.Println(self.weather.Rain)
			fmt.Println("watering = true")
		} else {
			watering = false
			self.config.MoistureThreshold = 0.5
			fmt.Println("watering = false")

		}
		if watering {
			self.runPump()

		}

		time.Sleep(2 * time.Second)
	}
}

func (self *Pi) ConnectTo(hub common.Hub) {
	//register self at hub
	hub.Register(self)
	go self.RoutineCheck()
}
