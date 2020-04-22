package weather

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"weather/internal/httpclient"
	"weather/internal/log"
)

type Client struct {
	baseURL string
	appID   string
	units   string
	lang    string
	d       httpclient.Doer
}

func NewClient(baseURL, apiKey, units, lang string, logger log.Logger) *Client {
	return &Client{
		baseURL: baseURL,
		units:   units,
		lang:    lang,
		appID:   apiKey,
		d:       httpclient.LoggingMiddleware(logger)(httpclient.DoerWrappedClient(5 * time.Second)),
	}
}

func (c *Client) GetByCity(q string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/weather/?q=%s&units=%s&lang=%s&appid=%s", c.baseURL, q, c.units, c.lang, c.appID), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.d.Do(req)
	if err != nil {
		return nil, err
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	_, body, err := allow2xxCodesOnly(overrideStatusCode(allow2xxCodesOnly(res.StatusCode, resBody, nil)))
	return body, err
}

func allow2xxCodesOnly(code int, body []byte, err error) (int, []byte, error) {
	if err != nil {
		return 0, nil, err
	}

	if code < 200 || code > 299 {
		return 0, nil, fmt.Errorf("http return code %d is not allowed", code)
	}

	return code, body, nil
}

func overrideStatusCode(code int, body []byte, err error) (int, []byte, error) {
	if err != nil {
		return 0, nil, err
	}

	var p InternalParams
	if err := json.Unmarshal(body, &p); err != nil {
		return 0, nil, err
	}

	if p.StatusCode != 0 && p.StatusCode != code {
		return p.StatusCode, body, nil
	}

	return code, body, nil
}

type Sys struct {
	Country string `json:"country"`
	ID      int64  `json:"id"`
	Sunrise int64  `json:"sunrise"`
	Sunset  int64  `json:"sunset"`
	Type    int64  `json:"type"`
}

type InternalParams struct {
	StatusCode int    `json:"cod"`
	CityID     int64  `json:"id"`
	TimeStamp  int64  `json:"dt"`
	Sys        Sys    `json:"sys"`
	Base       string `json:"base"`
}
