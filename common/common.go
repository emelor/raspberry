package common

type Pi interface {
	UpdateConfig(Configuration)
	UpdateWeather(Weather)
	//manual on	(minutes)
	//get status update

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
