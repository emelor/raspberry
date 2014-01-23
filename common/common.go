package common

import "time"

type Pi interface {
	UpdateConfig(Configuration)
	UpdateWeather(Weather)
	//status update
	//data update
}

type Hub interface {
	Register(Pi)

	//show webpage

}

type Weather struct {
	//rain during next 24 h
	Rain    float64
	RainNow float64
	//max temp during next 24 h

}

type Configuration struct {
	MoistureThreshold float64
	RainLimit         float64
	WeatherUrl        string
	ManualSetting     bool
	ManualUntil       time.Time
	//MinutesDaily		int
	//TimeDaily
}
