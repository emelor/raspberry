package pi

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"../common"
)

type Pi struct {
	//
	config        common.Configuration
	weather       common.Weather
	rain          float64
	wateringTimer *time.Timer
	pumpRunning   bool
	timeStamp     time.Time
}

func New() *Pi {
	return &Pi{}

}
func (self *Pi) saveSettings() {
	f, err := os.Create("settings.json")
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f)
	fmt.Println("saving:", self.config, "in", f.Name())
	err = enc.Encode(self.config)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func (self *Pi) saveWeather() {
	f, err := os.Create("weather.json")
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f)
	fmt.Println("saving:", self.weather, "in", f.Name())
	err = enc.Encode(self.weather)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}

func (self *Pi) getMoisture() float64 {

	moisture := rand.Float64()
	return moisture
}

func (self *Pi) UpdateConfig(config common.Configuration) {
	self.config = config
	//save new config to file
	self.saveSettings()

	if self.config.ManualOn {
		self.config.ManualOn = false
		self.config.ManualOff = false
		self.timeStamp = time.Now()
		fmt.Println("MinutesOn = ", self.config.MinutesOn)
		self.runPump(self.config.MinutesOn)
	}
	if self.config.ManualOff {
		self.config.ManualOn = false
		self.config.ManualOff = false
		self.timeStamp = time.Now()
		fmt.Println("MinutesOff = ", self.config.MinutesOff)
		self.stopPump(self.config.MinutesOff)
	}
}

func (self *Pi) UpdateWeather(weather common.Weather) {
	self.weather = weather
	fmt.Println("got weather", self.weather)
	//save new weather to file
	self.saveWeather()
}

func (self *Pi) runPump(minutes int) {
	//Start pump
	if self.pumpRunning {
		fmt.Println("Pump already running. New stop time in ", minutes, " minutes.")
	} else {
		fmt.Println("Pump starting. Stop time in ", minutes, " minutes.")
	}
	self.pumpRunning = true
	if self.wateringTimer != nil {
		self.wateringTimer.Stop()
	}
	self.wateringTimer = time.AfterFunc(time.Duration(minutes)*time.Minute, func() { self.stopPump(0) })

}

func (self *Pi) stopPump(minutes int) {
	//Stop pump, block watering for specified number of minutes
	//Start watering/evaluation only if timer has run out or if "ManualOn" is pressed in the UI
	self.config.ManualOn = false

	if self.pumpRunning {
		fmt.Println("Stopping pump")
	}
	self.pumpRunning = false
	if self.wateringTimer != nil {
		self.wateringTimer.Stop()
	}
}

func (self *Pi) RoutineCheck() {
	for {
		timeSince := time.Now().Sub(self.timeStamp)
		fmt.Println("timeSince = ", timeSince)
		minOn := time.Minute * time.Duration(self.config.MinutesOn)
		minOff := time.Minute * time.Duration(self.config.MinutesOff)
		//Evaluate only if enough time has elapsed since last override order (self.timeStamp)
		if timeSince > minOn && timeSince > minOff {
			var watering bool
			//Soil moisture sensor reading
			moisture := self.getMoisture()

			if (moisture < self.config.MoistureThreshold) && (self.weather.Rain <= self.config.RainLimit) {
				watering = true

				//fmt.Println(self.weather.Rain)

				fmt.Println("watering = true")
			} else {
				watering = false
				fmt.Println("watering = false")

			}
			if watering {
				fmt.Println("pump running (routine)")
				self.runPump(1)

			}
		}

		time.Sleep(10 * time.Second)
	}
}

func (self *Pi) ConnectTo(hub common.Hub) {
	//register self at hub
	hub.Register(self)
	go self.RoutineCheck()
}
