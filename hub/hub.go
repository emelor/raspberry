package hub

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"../common"
)

type Hub struct {
	config        common.Configuration
	weather       common.Weather
	rain          float64
	registeredPis []common.Pi
}

//Time, Precipitation and WeatherData structs are needed to mirror the structure of yr.no weather data
//...
type Time struct {
	From   string        `xml:"from,attr"`
	To     string        `xml:"to,attr"`
	Precip Precipitation `xml:"precipitation"`
}
type Precipitation struct {
	Value float64 `xml:"value,attr"`
}
type WeatherData struct {
	Times []Time `xml:"forecast>tabular>time"`
}

//...
//...

func New() *Hub {
	fmt.Println("new hub!")
	return &Hub{}
}
func (self *Hub) Register(pi common.Pi) {
	self.registeredPis = append(self.registeredPis, pi)
	pi.UpdateConfig(self.config)
	pi.UpdateWeather(self.weather)
}

func (self *Hub) transferSettings() {
	for _, pi := range self.registeredPis {
		pi.UpdateConfig(self.config)
	}
}

func (self *Hub) transferWeather() {
	for _, pi := range self.registeredPis {
		pi.UpdateWeather(self.weather)
	}
}

func (self *Hub) checkWeather() {
	for {
		client := new(http.Client)
		resp, err := client.Get(self.config.WeatherUrl)
		if err != nil {
			panic(err)
		}
		var weatherData WeatherData
		if err := xml.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
			panic(err)
		}
		var rainTotal = 0.0
		for i := 0; i < 8; i++ {
			rainMm := weatherData.Times[i].Precip.Value
			rainTotal = rainTotal + rainMm
			fmt.Println(rainMm)
			fmt.Println("Rain total: ")
			fmt.Println(rainTotal)
		}
		fmt.Println("***************************************************************")
		fmt.Println("fetching weather data from ", self.config.WeatherUrl)
		fmt.Println("***************************************************************")
		self.weather.Rain = rainTotal
		time.Sleep(time.Hour)
	}
}
