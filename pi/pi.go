package pi

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"../common"
)

//
//never reads settings from file, so why save to file?
//

type Pi struct {
	//
	config  common.Configuration
	weather common.Weather
	//rain          float64
	wateringTimer  *time.Timer
	pumpRunning    bool
	minutesWatered int
	moisture       float64
	logPath        string
}

func New() *Pi {
	//Default path for log files: current directory
	return &Pi{logPath: "."}

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
func init() {
	//Use current time as seed so that random function is more random
	rand.Seed(time.Now().UnixNano())
}

func (self *Pi) getMoisture() float64 {

	moisture := rand.Float64()
	self.moisture = moisture
	return moisture
}

func (self *Pi) UpdateConfig(config common.Configuration) {
	//Update config struct
	self.config = config
	//save new config to file
	self.saveSettings()
	if time.Now().Before(self.config.ManualUntil) {
		if self.config.ManualSetting {
			self.runPump(self.config.ManualUntil)
		} else {
			self.stopPump()
			fmt.Println("Pump stopped until ", self.config.ManualUntil)
		}
	}
}

func (self *Pi) UpdateWeather(weather common.Weather) {
	self.weather = weather
	fmt.Println("got weather", self.weather)
	//save new weather to file
	self.saveWeather()
}

func (self *Pi) runPump(until time.Time) {
	//Start pump
	if self.pumpRunning {
		fmt.Println("Pump already running. New stop time ", until.Format(time.Stamp))
	} else {
		fmt.Println("Pump starting. Stop time  ", until.Format(time.Stamp))
	}
	self.pumpRunning = true
	if self.wateringTimer != nil {
		self.wateringTimer.Stop()
	}
	self.wateringTimer = time.AfterFunc(until.Sub(time.Now()), func() { self.stopPump() })
}

func (self *Pi) stopPump() {
	//Stop pump, block watering for specified number of minutes
	//Start watering/evaluation only if timer has run out or if "ManualOn" is pressed in the UI

	if self.pumpRunning {
		fmt.Println("Stopping pump")
		self.pumpRunning = false
	}
	if self.wateringTimer != nil {
		self.wateringTimer.Stop()
	}
}

func (self *Pi) RoutineCheck() {
	for {
		if time.Now().After(self.config.ManualUntil) {
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
				self.minutesWatered += 1
				fmt.Println(self.minutesWatered)
				self.runPump(time.Now().Add(time.Minute))

			}
		}

		time.Sleep(10 * time.Second)
	}
}

func (self *Pi) ConnectTo(hub common.Hub) {
	//register self at hub
	hub.Register(self)
	go self.RoutineCheck()
	go self.histToFile()
}
