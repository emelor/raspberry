//Http handlers for the user interface: "update" (submit button), manual watering "on" and "off", respectively, and "show"

package hub

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
"common"
	"code.google.com/p/go.net/websocket"
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
	self.config.ManualSetting = true
	parsed, err := strconv.ParseInt(r.FormValue("minuteson"), 10, 64)
	if err != nil {
		panic(err)
	}
	self.config.ManualUntil = time.Now().Add(time.Minute * time.Duration(parsed))
	self.saveSettings()
	w.Header().Set("Location", "/")
	w.WriteHeader(303)
}

func (self *Hub) handleOff(w http.ResponseWriter, r *http.Request) {
	//Do everything that needs to be done before the UI is refreshed again
	fmt.Println("offoff")
	var err error
	self.config.ManualSetting = false
	parsed, err := strconv.ParseInt(r.FormValue("minutesoff"), 10, 64)
	if err != nil {
		panic(err)
	}
	//ManualUntil is set to make sure that automatic watering is not resumed until "minutesoff" minutes have passed
	self.config.ManualUntil = time.Now().Add(time.Minute * time.Duration(parsed))
	self.saveSettings()
	w.Header().Set("Location", "/")
	w.WriteHeader(303)
}

func (self *Hub) handleShow(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	fmt.Fprintf(w, `
	<html>
	<head>
	<script src="http://www.chartjs.org/docs/Chart.js"></script>
	</head>
	<body>
	<form action=/update>
Input lowest acceptable soil moisture level: <br/>
moisture (min: 0.0, max: 1.0): <input type='text' name='moisture' value='%v'> <br/><br/>


How many millimeters of rain in the forecast are required to delay watering? <br/>
Rain (mm): <input type='text' name='rainlimit' value='%v'> <br/><br/>


Input weather source URL for your location from <a href='http://fil.nrk.no/yr/viktigestader/verda.txt'>yr.no</a>: <br/>
URL: <input type='text' name='URL' value='%v' size=%v> <br/><br/> 

	
Total rain forecast in your location during the next 24 hours: %v mm </br><br/>

<input type='submit' value='Submit'>
</form>

<form action=/on>

<input type='submit' value='Start now'>
Start watering immediately. Watering will stop after the number of minutes input in the box below. <br/>
Number of minutes: <input type='text' name='minuteson' value='10'> <br/><br/>
</form>
<form action=/off>

<input type='submit' value='Stop now'>
Pause watering immediately. Automatic watering will be resumed after the number of minutes input in the box below. <br/>
Number of minutes: <input type='text' name='minutesoff' value='10'> <br/><br/>
</form>

<canvas id="myChart" width="1600" height="500"></canvas>
<script>

var ctx = document.getElementById("myChart").getContext("2d");
var data = %v;
var myNewChart = new Chart(ctx).Line(data);

</script>
	</body>
	</html>
	`, self.config.MoistureThreshold, self.config.RainLimit, self.config.WeatherUrl, len(self.config.WeatherUrl), self.weather.Rain, self.getDataPoints())

	//MinutesDaily		int
	//fmt.Fprintln(w, "Complement the automatic watering scheme by watering %v minutes a day, at %v <br/>", self.config.MinutesDaily self.config.TimeDaily)
	//fmt.Fprintf(w, "Number of minutes: <input type='text' name='minutesoff' value='%v'> <br/><br/>", self.config.MinutesOff)
	//TimeDaily

}

type JsonData struct {
	Data             []float64 `json:"data"`
	FillColor        string    `json:"fillColor"`
	StrokeColor      string    `json:"strokeColor"`
	PointColor       string    `json:"pointColor"`
	PointStrokeColor string    `json:"pointStrokeColor"`
}

type JsonDataSets struct {
	Labels   []interface{} `json:"labels"`
	Datasets []JsonData    `json:"datasets"`
}

/*
{
	labels: [1, 2, 3, 4, 5, 6, 7],
	datasets : [
		{
			data : [65,59,90,81,56,55,40]
		},
	]
}
fillColor : "rgba(220,220,220,0.5)",
			strokeColor : "rgba(220,220,220,1)",
			pointColor : "rgba(220,220,220,1)",
			pointStrokeColor : "#fff",
*/

func (self *Hub) getDataPoints() string {

	dataSets := JsonDataSets{
		Datasets: []JsonData{
			JsonData{
				//rain
				FillColor:        "rgba(255,170,255,0.5)",
				StrokeColor:      "rgba(255,170,255,1)",
				PointColor:       "rgba(255,170,255,1)",
				PointStrokeColor: "#faf",
			},
			JsonData{
				//moisture
				FillColor:        "rgba(170,255,255,0.5)",
				StrokeColor:      "rgba(170,255,255,1)",
				PointColor:       "rgba(170,255,255,1)",
				PointStrokeColor: "#aff",
			},
			JsonData{
				//minutes
				FillColor:        "rgba(255,255,170,0.5)",
				StrokeColor:      "rgba(255,255,170,1)",
				PointColor:       "rgba(255,255,170,1)",
				PointStrokeColor: "#ffa",
			},
		},
	}

	pi := self.registeredPis[0]
	from := time.Now().Add(time.Hour * -72)
	to := time.Now()
	dataStructs := pi.GetHistory(from, to)
	for _, point := range dataStructs {
		dataSets.Labels = append(dataSets.Labels, point.Time)
		dataSets.Datasets[0].Data = append(dataSets.Datasets[0].Data, point.Rain)
		dataSets.Datasets[1].Data = append(dataSets.Datasets[1].Data, point.Moisture)
		dataSets.Datasets[2].Data = append(dataSets.Datasets[2].Data, float64(point.Minutes))
	}
	blabla, err := json.Marshal(dataSets)
	if err != nil {
		panic(err)
	}
	return string(blabla)
}

func (self *Hub) handleWS(ws *websocket.Conn) {
	for {
		message := ""
		err := websocket.JSON.Receive(ws, &message)
		if err != nil {
			panic(err)
		}
		if message == "register"{
			self.Register(NewRemotePi(ws))
			
		}
}
type RemotePi struct{
	
}

func NewRemotePi(){
	
}
func (self *RemotePi) UpdateConfig(config common.Configuration){
	
}
func (self *RemotePi) UpdateWeather(weather common.Weather){
	
}
func (self *RemotePi) GetHistory(time.Time, time.Time)([]common.Data){
	
}

func (self *Hub) GetWSAddress() (address string) {
	address = "ws://127.0.0.1:25601/ws"
	return
}

func (self *Hub) serve() {
	http.HandleFunc("/update", http.HandlerFunc(self.handleUpdate))
	http.HandleFunc("/on", http.HandlerFunc(self.handleOn))
	http.HandleFunc("/off", http.HandlerFunc(self.handleOff))
	http.HandleFunc("/", http.HandlerFunc(self.handleShow))
	http.Handle("/ws", websocket.Handler(self.handleWS))
	go func() {
		if err := http.ListenAndServe(":25601", nil); err != nil {
			panic(err)
		}
	}()

}
