package main

/*
Webserver implemented on RaspberryPi hardware
Collect weather data from a website/api (or another server masquerading as a sensor)
Moisture sensor gives soil moisture data
RaspberryPi decides whether or not to turn on watering system/valve
Web server lets user configure the system through web form

*/

import (
	//"encoding/json"
	"fmt"
	//"sync"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"encoding/xml"
)

var minMoisture = 0.5
var rain = false

//Web server:
func handle(w http.ResponseWriter, r *http.Request) {
	willRain()
	fmt.Println(r)
	r.ParseForm()
	fmt.Println(r.Form["moisture"])
	w.Header().Set("Content-Type", "text/html")
	if len(r.Form["moisture"]) > 0 {
		var err error
		minMoisture, err = strconv.ParseFloat(r.Form["moisture"][0], 64)
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}

	if len(r.Form["rain"]) > 0 {
		rain = r.Form["rain"][0] == "1"
		

	}

	fmt.Fprintln(w, "<form>")
	fmt.Fprintln(w, "Input lowest acceptable soil moisture level: <br/>")
	fmt.Fprintf(w, "moisture: <input type='text' name='moisture' value='%v'> <br/>", minMoisture)
	fmt.Fprintln(w, "Will it rain at any time during the next 24 hours?: <br/>")
	if rain {
		fmt.Fprintln(w, "<input type='radio' checked='checked' name='rain' value=1> Yes <br/>")
		fmt.Fprintln(w, "<input type='radio' name='rain' value=0> No")
	} else {
		fmt.Fprintln(w, "<input type='radio' name='rain' value=1> Yes <br/>")
		fmt.Fprintln(w, "<input type='radio' checked='checked' name='rain' value=0> No")
	}
	fmt.Fprintln(w, "<input type='submit' value='Submit'>")
	fmt.Fprintln(w, "</form>")
}

func getMoisture() float64 {
	moisture := rand.Float64()
	return moisture

}

type Precipitation struct {
	Value float64 `xml:"value,attr"`
}

type WeatherData struct {
	Precipitations []Precipitation `xml:"forecast>tabular>time>precipitation"`
}

func willRain() float64 {
	client := new(http.Client)
	resp, err := client.Get("http://www.yr.no/stad/Andorra/Encamp/Vila/varsel.xml")
	if err != nil {
		panic(err)
	}
	var weatherData WeatherData
	if err := xml.NewDecoder(resp.Body).Decode(&weatherData); err != nil {
		panic(err)
	}
	fmt.Println(weatherData)
	return 0
}

func runPump() {
	fmt.Println("Pump is running")
	fmt.Println(rain)
}

func checkLoop() {
	for {
		var watering bool
		moisture := getMoisture()
		
		if (moisture < minMoisture) && !(rain) {
			watering = true
		} else {
			watering = false
		}
		if watering {
			runPump()
		} else {
			fmt.Println("Pump is not running")
			fmt.Println(rain)
		}
		time.Sleep(2 * time.Second)

	}
}

func main() {
	go checkLoop()
	http.HandleFunc("/", http.HandlerFunc(handle))
	if err := http.ListenAndServe(":25601", nil); err != nil {
		panic(err)
	}
}
