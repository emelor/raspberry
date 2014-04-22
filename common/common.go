package common

import "time"

type Pi interface {
	UpdateConfig(Configuration)
	UpdateWeather(Weather)
	GetHistory(time.Time, time.Time) []Data
}

type Hub interface {
	Register(Pi)
}

type Weather struct {
	//rain during next 24 h
	Rain float64
	//rain during next 3 hours (time resolution of yr.no weather data)
	RainNow float64
}

type Message struct {
	FunctionName string
	ConfigBody   *Configuration
	WeatherBody  *Weather
	HistoryBody  *HistoryBody
	Data         []Data
	MessageID    int64
}

type HistoryBody struct {
	From time.Time
	To   time.Time
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

type Data struct {
	Rain     float64
	Minutes  int
	Moisture float64
	Time     time.Time
}
