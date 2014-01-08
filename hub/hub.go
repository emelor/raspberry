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
	rain          float64
	registeredPis []common.Pi
	weather       common.Weather
}

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

//func New() (returnValue *Hub)
func New() *Hub {
	fmt.Println("new hub!")
	return &Hub{}
}
func (self *Hub) Register(pi common.Pi) {
	self.registeredPis = append(self.registeredPis, pi)
}

func (self *Hub) transferSettings() {
	for _, pi := range self.registeredPis {
		pi.UpdateConfig(self.config)
	}
}

func (self *Hub) transferWeather() {
	for i := 0; i < len(self.registeredPis); i++ {
		self.registeredPis[i].Weather = self.weather
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
