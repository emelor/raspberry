package common

type Pi interface {
	//NewSettings
	//manual on	(minutes)
	//get status update
	//routine (water?, wait, water?, wait)

}

type Hub interface {
	Register(Pi)
	//get weather
	//request data update
	//show webpage
	//push settings
	//push weather
	//register new node

}

type Weather struct {
	//rain during next 24 h
	Rain float64
	//max temp during next 24 h

}

type Configuration struct {
	MoistureThreshold float64
	RainLimit         float64
	WeatherUrl        string
}
