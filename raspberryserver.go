package main

/*

Webserver implemented on RaspberryPi hardware
Collect weather data from a website/api
Moisture sensor gives soil moisture data
RaspberryPi decides whether or not to turn on watering system/valve
Web server lets user configure the system through web form

*/

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"encoding/xml"
)

var minMoisture = 0.5
var rainLimit = 4.0
var rain = false
var rainStruct = &rainCheck{}
var weatherUrl = "http://www.yr.no/stad/Sverige/Stockholm/Stockholm/varsel.xml"
type configuration struct{
	moist float64
	rain float64
	url string
}


//Web server:
func handle(w http.ResponseWriter, r *http.Request) {
	rainStruct.willRain()
	if rainStruct.rainTotal>0 {
		rain = true
		fmt.Println("rain is true")
	}
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

	if len(r.Form["rainlimit"]) > 0 {
		var err error
		rainLimit, err = strconv.ParseFloat(r.Form["rainlimit"][0], 64)
		if err != nil {
			fmt.Fprintln(w, err)
		}
	}


	if len(r.Form["URL"]) > 0 {
		if weatherUrl != r.Form["URL"][0]{
			rainStruct.previousTime = time.Now()
		}
		weatherUrl = r.Form["URL"][0]
		}
	fmt.Fprintln(w, "<form>")
	fmt.Fprintln(w, "Input lowest acceptable soil moisture level: <br/>")
	fmt.Fprintf(w, "moisture (min: 0.0, max: 1.0): <input type='text' name='moisture' value='%v'> <br/><br/>", minMoisture)

	fmt.Fprintln(w, "How many millimeters of rain in the forecast are required to delay watering? <br/>")
	fmt.Fprintf(w, "Rain (mm): <input type='text' name='rainlimit' value='%v'> <br/><br/>", rainLimit)

	fmt.Fprintln(w, "Input weather source URL for your location from <a href='http://fil.nrk.no/yr/viktigestader/verda.txt'>yr.no</a>: <br/>")
	fmt.Fprintf(w, "URL: <input type='text' name='URL' value='%v' size=%v> <br/><br/>", weatherUrl, len(weatherUrl))

	fmt.Fprintln(w, "Total rain forecast in your location during the next 24 hours: ")
	fmt.Fprintln(w, rainStruct.rainTotal)
	fmt.Fprintln(w, " mm </br><br/>" )
	//put in override "water anyway? buttons: today, always, no"
	
	fmt.Fprintln(w, "<input type='submit' value='Submit'>")
	fmt.Fprintln(w, "</form>")
}

func getMoisture() float64 {
	moisture := rand.Float64()
	return moisture

}

type rainCheck struct {
	previousTime time.Time 
	rainTotal float64
}

type Time struct {
	From string `xml:"from,attr"`
	To string `xml:"to,attr"`
	Precip Precipitation `xml:"precipitation"`
}
type Precipitation struct {
	Value float64 `xml:"value,attr"`
}

type WeatherData struct {
	Times []Time `xml:"forecast>tabular>time"`
}

func (self *rainCheck) willRain() {
	if time.Now().Sub(self.previousTime)<time.Hour{
		fmt.Println("willRain func fresh timestamp")
	}else{
	client := new(http.Client)
	resp, err := client.Get(weatherUrl)
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
		fmt.Println("fetching weather data from ", weatherUrl)
		fmt.Println("***************************************************************")
		self.rainTotal = rainTotal
		self.previousTime = time.Now()
	}
	//fmt.Println(weatherData)
}

func runPump() {
	fmt.Println("Pump is running")
	//fmt.Println(rain)
}

func checkLoop() {
	for {
		var watering bool
		moisture := getMoisture()
		
		if (moisture < minMoisture) && (rainStruct.rainTotal < rainLimit) {
			watering = true
		} else {
			watering = false
		}
		if watering {
			runPump()
		} else {
			fmt.Println("Pump is not running")
			fmt.Println(minMoisture)
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
