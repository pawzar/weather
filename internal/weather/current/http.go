package current

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"weather/internal/log"
	"weather/internal/weather"
)

func NewHttpService(c *weather.Client) *HttpService {
	return &HttpService{c: c}
}

type HttpService struct {
	c *weather.Client
}

func (s *HttpService) ByCity(wq *WeatherQuery) (*Weather, error) {
	q := wq.City
	if wq.Country != "" {
		q = q + "," + wq.Country
	}

	jp, err := s.c.GetByCity(q)
	if err != nil {
		return nil, err
	}

	var p byCity
	if err := json.Unmarshal(jp, &p); err != nil {
		return nil, err
	}

	t := time.Unix(p.UTCInSeconds, 0)

	return &Weather{
		CalculationTime:      t,
		City:                 p.Name,
		Location:             p.Location,
		CloudinessPercentage: int(p.Clouds.All),
		Wind:                 p.Wind,
		Conditions:           p.Conditions,
		VisibilityMetres:     int(p.Visibility),
		Descriptions:         p.Descriptions,
	}, nil
}

func MultipleCitySearch(s Service, logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Debugf("< %s %s %s", r.Method, r.URL, r.Proto)

		wqs, err := parseQueryString(r)
		if err != nil {
			logger.Debugf(" invalid input data")
			w.WriteHeader(400)

			return
		}

		count := callAPIAndWriteData(wqs, s, logger, w)

		logger.Debugf("> bytes sent: %d", count)
	}
}

func callAPIAndWriteData(wqs []*WeatherQuery, s Service, logger log.Logger, w http.ResponseWriter) int {
	weatherChannel := make(chan *Weather, 3)
	var wg1 sync.WaitGroup
	for _, wq := range wqs {
		wg1.Add(1)

		go func(wq *WeatherQuery) {
			defer wg1.Done()

			weatherInCity, err := s.ByCity(wq)

			if err != nil {
				logger.Warningf("failed to get weather data: %s", err)
				return
			}

			weatherChannel <- weatherInCity
		}(wq)
	}

	var sum int

	var wg2 sync.WaitGroup
	wg2.Add(1)
	go func() {
		defer wg2.Done()

		for weatherInCity := range weatherChannel {
			n, err := fmt.Fprintf(w, "%+v\n", weatherInCity)
			if err != nil {
				logger.Errorf("failed to write: %s", err)

				return
			}
			sum = sum + n
		}
	}()

	wg1.Wait()
	close(weatherChannel)

	wg2.Wait()
	return sum
}

func parseQueryString(r *http.Request) ([]*WeatherQuery, error) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}

	var wqs []*WeatherQuery

	if cities, exists := params["c"]; !exists {
		return nil, errors.New("no 'city' given in query")
	} else {
		for _, c := range cities {
			wqs = append(wqs, &WeatherQuery{
				City: c,
			})
		}
	}

	if countryCodes, exists := params["cc"]; !exists {
		return wqs, nil
	} else {
		for k, cc := range countryCodes {
			if k == len(wqs) {
				break
			}

			wqs[k].Country = cc
		}
	}

	return wqs, err
}
