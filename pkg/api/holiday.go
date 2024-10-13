package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		ct.Time = time.Time{}
		return
	}
	ct.Time, err = time.Parse(time.DateOnly, s)
	return
}

type Holiday struct {
	Date        CustomTime `json:"date"`
	LocalName   string     `json:"localName"`
	Name        string     `json:"name"`
	CountryCode string     `json:"countryCode"`
	Fixed       bool       `json:"fixed"`
	Global      bool       `json:"global"`
	Counties    *[]string  `json:"counties"`
	LaunchYear  *int       `json:"launchYear"`
	Types       []string   `json:"types"`
}

type comboKey struct {
	date    int
	country string
}
type HolidayService struct {
	cache map[comboKey][]*Holiday
}

func NewHoliday() *HolidayService {
	return &HolidayService{cache: make(map[comboKey][]*Holiday, 0)}
}

func (h HolidayService) findHoliday(date time.Time) bool { //TODO Split into two - holiday and day before holiday
	b, ok := h.cache[comboKey{date.Year(), "DE"}]
	if !ok {
		return false
	}
	_, ok = slices.BinarySearchFunc(b, date.UTC(), func(holiday *Holiday, d time.Time) int {
		return holiday.Date.Compare(d)
	})
	if ok {
		return ok
	}
	dayBefore := date.Add(24 * time.Hour)

	_, ok = slices.BinarySearchFunc(b, dayBefore.UTC(), func(holiday *Holiday, d time.Time) int {
		return holiday.Date.Compare(d)
	})

	return ok
}

func (h HolidayService) getHoliday(date time.Time, country string) {
	resp, err := http.Get(fmt.Sprintf("https://date.nager.at/api/v3/publicholidays/%d/%s", date.Year(), country))
	if err != nil {
		log.Error(err)
		return
	}
	var holidays []*Holiday
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = json.Unmarshal(body, &holidays)
	if err != nil {
		log.Error(err)
		return
	}
	slices.SortFunc(holidays, func(a, b *Holiday) int {
		return a.Date.Compare(b.Date.Time)
	})
	h.cache[comboKey{date.Year(), country}] = holidays

}
func (h HolidayService) IsPublicHoliday(date time.Time, country string) bool {
	_, ok := h.cache[comboKey{date.Year(), country}]
	if ok {
		return h.findHoliday(date)
	} else {
		h.getHoliday(date, country)
		return h.findHoliday(date)
	}

}
