//Save watering and moisture data to file. Json encoding.
//A new file every day, named by date, makes it easy to retrieve data from a given interval
//History is saved on the Raspberry Pi, and sent to browser via the hub on request

package pi

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../common"
)

func (self *Pi) GetHistory(from, to time.Time) []common.Data {
	dataPoints := []common.Data{}
	//open current directory
	dir, err := os.Open(self.logPath)
	if err != nil {
		panic(err)
	}
	//Get file names in slice
	names, err := dir.Readdirnames(-1)
	if err != nil {
		panic(err)
	}
	for i, name := range names {
		name = strings.TrimRight(name, `.json`)
		//Check if file name is a date
		parsedTime, err := time.Parse("2006January02", name)
		if err == nil {
			//Is file date stamp later than, or equal to, 'from' AND earlier than, or equal to, 'to'?
			if (parsedTime.Equal(from) || parsedTime.After(from)) && (parsedTime.Equal(to) || parsedTime.Before(to)) {
				//for matching files, go through and append structs to dataPoints
				nameString := names[i]
				file, err := os.Open(filepath.Join(self.logPath, nameString))
				if err != nil {
					panic(err)
				}

				dec := json.NewDecoder(file)
				for {
					var newDataPoint common.Data
					err := dec.Decode(&newDataPoint)
					if err != nil {
						fmt.Println("End of file", err)
						break
					}
					dataPoints = append(dataPoints, newDataPoint)

				}
			}

		}

	}

	return dataPoints
}

func (self *Pi) fetchData() common.Data {
	var d common.Data
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
		time.Sleep(time.Minute)
	}
}
