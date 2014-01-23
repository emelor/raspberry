//Save watering and moisture data to file. Json encoding.
//A new file every day, named by date, makes it easy to retrieve data from a given interval

package pi

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type data struct {
	Rain     float64
	Minutes  int
	Moisture float64
	Time     time.Time
}

func (self *Pi) fetchData() data {
	var d data
	d.Rain = self.weather.RainNow
	d.Minutes = self.minutesWatered
	d.Moisture = self.moisture
	d.Time = time.Now()
	self.minutesWatered = 0
	return d
}

func (self *Pi) histToFile() {
	for {
		d := self.fetchData()
		year, month, day := time.Now().Date()
		filename := fmt.Sprintf("%v%v%v.json", year, month, day)
		//Open file by 'filename' in append mode. If no such file  exists, create it first.
		//0666 means that file will be created readable and writeable by all
		file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		//file, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(file)
		fmt.Println("saving:", d, "in", file.Name())
		err = enc.Encode(d)
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}

		//encode and append data
		time.Sleep(time.Hour)
	}
}
