package common

type Pi interface{
//NewSettings
//manual on	(minutes)
//get status update
//routine (water?, wait, water?, wait)
	
}

type Hub interface{
//get weather	
//request data update
//show webpage
//push settings
//push weather
//register new node

}

type Configuration struct{
	MoistureThreshold float64
	RainLimit float64
	WeatherUrl string
}