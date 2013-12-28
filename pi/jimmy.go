package main

import "fmt"

func main() {
	fmt.Println("Hello, playground")
	
	chan force_stop
	chan force_start
	
	struct settings{
		timer Date.
	}
	
	open_settings()
	
	go func() {
	start_direkt = false
	for {
		if !(start_direkt) {
			fetch_weather()
			read_soil_moisture()
			des, mins := make_decision(settings)
		}
		if des == "water" or start_direkt {
		        start_direkt = false
			start_pump()
			select {
			  case time.Delay(mins):
			    break
			  case <- force_stop
			   break
			}
			
			stop_pump()
		}
		
		select {
		 case Time.Delay("1h"):
			break
		 case <- force_start
		      start_direkt = true
			break
		}
			
	
	}
	}()
	
	ws.listen(func() {
	swith(command) {
	  case "Vattna":
	        force_start <- 1
	
	  case "Sluta Vattna":
		force_stop <- 1
	
	   case "SetSetting":
	     settings[timer] = 5min
	
	  case "Get Log":
	})
}
