package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type PrayerTimeAPI struct {
	client *resty.Client
	url    string
}

func NewPrayerTimeAPI() *PrayerTimeAPI {
	return &PrayerTimeAPI{
		client: resty.New(),
		url:    "http://api.aladhan.com/v1/calendar/{year}/{month}",
	}
}

func (p *PrayerTimeAPI) GetYearData(year int, latitude, longitude float64) MidnightTimes {
	ch := make(chan MidnightTimes, 12)
	for i := 1; i <= 12; i++ {
		go func() {
			ch <- p.GetMonthData(year, i, latitude, longitude)
		}()
	}

	times := make(MidnightTimes)
	for i := 1; i <= 12; i++ {
		monthTimes := <-ch
		for k, v := range monthTimes {
			times[k] = v
		}
	}

	return times
}

func (p *PrayerTimeAPI) GetMonthData(
	year, month int,
	latitude, longitude float64,
) MidnightTimes {
	data, err := p.getRawData(year, month, latitude, longitude)
	if err != nil {
		return nil
	}

	times, err := p.parseRawData(data)
	if err != nil {
		return nil
	}

	return times
}

func (p *PrayerTimeAPI) getDate(respData RespData) string {
	return strings.ReplaceAll(respData.Date.Gregorian.Date, "-", ".")
}

func (p *PrayerTimeAPI) getMidnight(respData RespData) string {
	return respData.Timings.Midnight[:5]
}

func (p *PrayerTimeAPI) parseRawData(data []byte) (MidnightTimes, error) {
	var unmarshalled Resp
	err := json.Unmarshal(data, &unmarshalled)
	if err != nil {
		return nil, err
	}

	times := make(MidnightTimes)
	for _, respData := range unmarshalled.Data {
		times[p.getDate(respData)] = p.getMidnight(respData)
	}

	return times, nil
}

func (p *PrayerTimeAPI) getRawData(
	year, month int,
	latitude, longitude float64,
) ([]byte, error) {
	resp, err := p.client.R().SetPathParams(
		map[string]string{
			"year":  strconv.Itoa(year),
			"month": strconv.Itoa(month),
		},
	).SetQueryParams(map[string]string{
		"latitude":  fmt.Sprintf("%f", latitude),
		"longitude": fmt.Sprintf("%f", longitude),
		// standard params
		"method":       "3", // Muslim World League
		"midnightMode": "1", // From sunset to fajr
	}).Get(p.url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("status code: %d, resp body: %s", resp.StatusCode(), resp.String())
	}

	return resp.Body(), nil
}
