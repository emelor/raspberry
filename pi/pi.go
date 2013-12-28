package pi

/*

Implemented on RaspberryPi hardware
Moisture sensor gives soil moisture data
RaspberryPi decides whether or not to turn on watering system/valve
Actuation loop takes in:
 - weather data (via server, mm of rain in the next 24 hours) 
 - user settings (from server, via saved file on RPi)
 - moisture sensor reading (from i/o pins (initially mockup))


*/

import (
	"fmt"
	"math/rand"
	"time"
	"encoding/json"
	"os"
)



var rainStruct = &rainCheck{}

type rainCheck struct {
	//Time stamp for last weather update, rain forecast at that time
	previousTime time.Time 
	rainTotal float64
}


var configStruct = &configuration{
	moistureThreshold: 0.5,
	RainLimit: 2,
}
/*type configuration struct{
	//User settings: soil moisture and rain forecast thresholds
	moistureThreshold float64
	rainLimit float64
	//insert other parameters
}
*/

func soilMoisture() float64 {
	//Soil moisture placeholder
	moisture := rand.Float64()
	return moisture

}

func (self *configuration) saveSettings() {
	//Save latest settings to file
	f, err := os.Create("settings.json")
	if err != nil {
		panic(err)
	}
	enc := json.NewEncoder(f)
	fmt.Println("saving:", self, "in", f.Name())
	err = enc.Encode(self)
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
}


func runPump() {
	//Actuator placeholder
	pumpRunning = true
	fmt.Println("Pump is running")

}


func RoutineCheck() {
	for {
		var watering bool
		//Soil moisture sensor reading
		moisture := getMoisture()

		if (moisture < configStruct.moistureThreshold) && (rainStruct.rainTotal < configStruct.RainLimit) {
			watering = true
		} else {
			watering = false
		}
		if watering {
			runPump()
			
		} 
		
	time.Sleep(2 * time.Second)
	}
}

func (self *Pi) ConnectTo(hub common.Hub) {
	//register self at hub
	hub.Register(self)		
}

func (self *rainCheck) willRain() {
	if time.Now().Sub(self.previousTime)<time.Hour{
		fmt.Println("willRain func fresh timestamp")
	}else{
		self.previousTime = time.Now()
		client := new(http.Client)
		resp, err := client.Get(configStruct.Url)
		if err != nil {
			panic(err)
		}
		var weatherData WeatherData
		if err := xml.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
			panic(err)
		}
		var rainTotal = 0.0
		for i:=0; i<8; i++ {
			rainMm := weatherData.Times[i].Precip.Value
			rainTotal = rainTotal + rainMm
			fmt.Println(rainMm)
			fmt.Println("Rain total: ")
			fmt.Println(rainTotal)
		}
		fmt.Println("***************************************************************")
		fmt.Println("fetching weather data from ", configStruct.Url)
		fmt.Println("***************************************************************")
		self.rainTotal = rainTotal
	}
	//fmt.Println(weatherData)
}



/*func main() {


	f, err := os.Open("settings.json")
	//Import latest settings from file
	if err == nil{
		json.NewDecoder(f).Decode(configStruct)
		fmt.Println(configStruct)
		f.Close()
	}else{
		fmt.Println(err)
		//will actuate with default values set in the beginning of the file 
	}
	//If new user settings from server:
	//If actuation manual override (minutes) -> run pump(minutes)
	//If set time && extra watering(minutes) -> run pump(minutes)
	go checkLoop()
	//
	//	http.HandleFunc("/", http.HandlerFunc(handle))
	//	if err := http.ListenAndServe(":25601", nil); err != nil {
		//		panic(err)
	//	}
}
*/