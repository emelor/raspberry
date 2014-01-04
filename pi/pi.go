package pi

import (
	"math/rand"
	"time"

	"../common"
)

type Pi struct {
	//
	config common.Configuration
	rain   float64
}

func New() *Pi {
	return &Pi{}

}

func (self *Pi) getMoisture() float64 {

	moisture := rand.Float64()
	return moisture
}

func (self *Pi) runPump() {

}

func (self *Pi) RoutineCheck() {
	for {
		var watering bool
		//Soil moisture sensor reading
		moisture := self.getMoisture()

		if (moisture < self.config.MoistureThreshold) && (self.rain < self.config.RainLimit) {
			watering = true
		} else {
			watering = false
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
