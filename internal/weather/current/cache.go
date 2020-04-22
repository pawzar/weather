package current

import (
	"fmt"
	"sync"
	"time"
)

func NewCachingWrapper(s Service, c Cache) ServiceWrapper {
	return func(q *WeatherQuery) (*Weather, error) {
		if w, found := c.Fetch(*q); found {
			return &w, nil
		}

		w, err := s.ByCity(q)
		if err != nil {
			return nil, err
		}

		c.Store(*q, *w)

		return w, nil
	}
}

type Cache interface {
	Fetch(WeatherQuery) (Weather, bool)
	Store(WeatherQuery, Weather)
}

type entry struct {
	expiration time.Time
	data       *Weather
}
type InMemoryCache map[string]entry

var lock = sync.RWMutex{}

func (c InMemoryCache) Fetch(q WeatherQuery) (Weather, bool) {
	lock.RLock()
	defer lock.RUnlock()

	key := fmt.Sprintf("%s-%s", q.City, q.Country)

	if e, present := c[key]; present {
		//log.Printf("fetched: [%s] %s: %v", key, e.expiration, e.data)

		if time.Now().Before(e.expiration) {
			return *e.data, true
		}
	}

	return Weather{}, false
}

func (c InMemoryCache) Store(q WeatherQuery, w Weather) {
	lock.Lock()
	defer lock.Unlock()

	key := fmt.Sprintf("%s-%s", q.City, q.Country)
	exp := time.Now().Truncate(time.Minute).Add(time.Minute)

	//log.Printf("storing [%s] %s: %v", key, exp, w)

	c[key] = entry{
		expiration: exp,
		data:       &w,
	}
}
