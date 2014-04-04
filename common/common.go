package common

import "time"

type Pi interface {
	UpdateConfig(Configuration)
	UpdateWeather(Weather)
	GetHistory(time.Time, time.Time) []Data
	//status update
	//data update
}

type Hub interface {
	Register(Pi)
	//UpdateHistory(History)

	//show webpage

}

type Weather struct {
	//rain during next 24 h
	Rain    float64
	RainNow float64
	//max temp during next 24 h

}
type Message struct {
	FunctionName string
	ConfigBody   *Configuration
	WeatherBody  *Weather
	HistoryBody  *HistoryBody
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
