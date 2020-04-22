package current

import (
	"time"
)

type Service interface {
	ByCity(wq *WeatherQuery) (*Weather, error)
}

type ServiceWrapper func(wq *WeatherQuery) (*Weather, error)

func (x ServiceWrapper) ByCity(wq *WeatherQuery) (*Weather, error) {
	return x(wq)
}

type byCity struct {
	Name               string `json:"name"`
	TimeShiftInSeconds int64  `json:"timezone"`
	UTCInSeconds       int64  `json:"dt"`
	Clouds             struct {
		All int64 `json:"all"`
	} `json:"clouds"`
	Location     Location      `json:"coord"`
	Conditions   Conditions    `json:"main"`
	Visibility   int64         `json:"visibility"`
	Descriptions []Description `json:"weather"`
	Wind         Wind          `json:"wind"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Wind struct {
	Deg   int64   `json:"deg"`
	Speed float64 `json:"speed"`
}

type Conditions struct {
	FeelsLike float64 `json:"feels_like"`
	Humidity  int64   `json:"humidity"`
	Pressure  int64   `json:"pressure"`
	Temp      float64 `json:"temp"`
	TempMax   float64 `json:"temp_max"`
	TempMin   float64 `json:"temp_min"`
}

type Description struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
	ID          int64  `json:"id"`
	Main        string `json:"main"`
}

type Weather struct {
	CalculationTime      time.Time
	City                 string
	Location             Location
	CloudinessPercentage int
	Wind                 Wind
	Conditions           Conditions
	VisibilityMetres     int
	Descriptions         []Description
}

type WeatherQuery struct {
	City    string
	Country string
}
