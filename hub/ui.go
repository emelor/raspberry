package hub

import (
	"fmt"
	"net/http"
	"strconv"
)

//handleUpdate will receive parameters from the browser, update the settings, and redirect to / (handleShow)
func (self *Hub) handleUpdate(w http.ResponseWriter, r *http.Request) {
	//Do everything that needs to be done before the UI is refreshed again
	var err error
	self.config.MoistureThreshold, err = strconv.ParseFloat(r.FormValue("moisture"), 64)
	if err != nil {
		panic(err)
	}
	self.config.RainLimit, err = strconv.ParseFloat(r.FormValue("rainlimit"), 64)
	if err != nil {
		panic(err)
	}
	//self.config.ManualOn, err = strconv.ParseBool(r.FormValue("manualon"))
	//if err != nil {
	//	panic(err)
	//}
	self.config.WeatherUrl = r.FormValue("URL")
	self.saveSettings()
	fmt.Println("...", r.FormValue("URL"))
	w.Header().Set("Location", "/")
	w.WriteHeader(303)
}

func (self *Hub) handleOn(w http.ResponseWriter, r *http.Request) {
	//Do everything that needs to be done before the UI is refreshed again
	self.config.ManualOn = true
	parsed, err := strconv.ParseInt(r.FormValue("minuteson"), 10, 64)
	if err != nil {
		panic(err)
	}
	self.config.MinutesOn = int(parsed)
	self.saveSettings()
	w.Header().Set("Location", "/")
	w.WriteHeader(303)
}

func (self *Hub) handleOff(w http.ResponseWriter, r *http.Request) {
	//Do everything that needs to be done before the UI is refreshed again
	fmt.Println("offoff")
	var err error
	self.config.ManualOff = true
	parsed, err := strconv.ParseInt(r.FormValue("minutesoff"), 10, 64)
	if err != nil {
		panic(err)
	}
	self.config.MinutesOff = int(parsed)
	self.saveSettings()
	w.Header().Set("Location", "/")
	w.WriteHeader(303)
}

func (self *Hub) handleShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(w, "<form action=/update>")
	//config.MoistureThreshold
	fmt.Fprintln(w, "Input lowest acceptable soil moisture level: <br/>")
	fmt.Fprintf(w, "moisture (min: 0.0, max: 1.0): <input type='text' name='moisture' value='%v'> <br/><br/>", self.config.MoistureThreshold)

	//config.RainLimit
	fmt.Fprintln(w, "How many millimeters of rain in the forecast are required to delay watering? <br/>")
	fmt.Fprintf(w, "Rain (mm): <input type='text' name='rainlimit' value='%v'> <br/><br/>", self.config.RainLimit)

	//config.WeatherUrl
	fmt.Fprintln(w, "Input weather source URL for your location from <a href='http://fil.nrk.no/yr/viktigestader/verda.txt'>yr.no</a>: <br/>")
	fmt.Fprintf(w, "URL: <input type='text' name='URL' value='%v' size=%v> <br/><br/>", self.config.WeatherUrl, len(self.config.WeatherUrl))

	//weather.Rain
	fmt.Fprintln(w, "Total rain forecast in your location during the next 24 hours: ")
	fmt.Fprintln(w, self.weather.Rain)
	fmt.Fprintln(w, " mm </br><br/>")

	//Submit button
	fmt.Fprintln(w, "<input type='submit' value='Submit'>")
	fmt.Fprintln(w, "</form>")

	fmt.Fprintln(w, "<form action=/on>")
	//ManualOn          bool
	fmt.Fprintf(w, "<input type='submit' value='Start now'>")
	//MinutesOn         int
	fmt.Fprintln(w, "Start watering immediately. Watering will stop after the number of minutes input in the box below. <br/>")
	fmt.Fprintf(w, "Number of minutes: <input type='text' name='minuteson' value='%v'> <br/><br/>", self.config.MinutesOn)
	fmt.Fprintln(w, "</form>")

	fmt.Fprintln(w, "<form action=/off>")
	//ManualOff         bool
	//MinutesOff		int
	fmt.Fprintln(w, "<input type='submit' value='Stop now'>")
	fmt.Fprintln(w, "Pause watering immediately. Automatic watering will be resumed after the number of minutes input in the box below. <br/>")
	fmt.Fprintf(w, "Number of minutes: <input type='text' name='minutesoff' value='%v'> <br/><br/>", self.config.MinutesOff)
	fmt.Fprintln(w, "</form>")
	//MinutesDaily		int
	//fmt.Fprintln(w, "Complement the automatic watering scheme by watering %v minutes a day, at %v <br/>", self.config.MinutesDaily self.config.TimeDaily)
	//fmt.Fprintf(w, "Number of minutes: <input type='text' name='minutesoff' value='%v'> <br/><br/>", self.config.MinutesOff)
	//TimeDaily

}

func (self *Hub) serve() {
	http.HandleFunc("/update", http.HandlerFunc(self.handleUpdate))
	http.HandleFunc("/on", http.HandlerFunc(self.handleOn))
	http.HandleFunc("/off", http.HandlerFunc(self.handleOff))
	http.HandleFunc("/", http.HandlerFunc(self.handleShow))
	go func() {
		if err := http.ListenAndServe(":25601", nil); err != nil {
			panic(err)
		}
	}()

}
