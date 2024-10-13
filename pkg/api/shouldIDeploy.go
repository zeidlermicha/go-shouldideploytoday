package api

import (
	"os"
	"time"

	"golang.org/x/exp/rand"
	"gopkg.in/yaml.v3"
)

type ShoudldIDeploy struct {
	reasons        *Reasons
	holidayService *HolidayService
}

type Reasons struct {
	ToDeploy           []string `yaml:"toDeploy"`
	ToNotDdeploy       []string `yaml:"toNotDdeploy"`
	ThursdayAfternoon  []string `yaml:"thursdayAfternoon"`
	FridayAfternoon    []string `yaml:"fridayAfternoon"`
	Friday13Th         []string `yaml:"friday13th"`
	Afternoon          []string `yaml:"afternoon"`
	Weekend            []string `yaml:"weekend"`
	DayBeforeChristmas []string `yaml:"dayBeforeChristmas"`
	Christmas          []string `yaml:"christmas"`
	NewYear            []string `yaml:"newYear"`
}

func NewShouldIDeploy(path string) (*ShoudldIDeploy, error) {
	reasons, err := loadReasons(path)
	if err != nil {
		return nil, err
	}

	return &ShoudldIDeploy{reasons: reasons, holidayService: NewHoliday()}, nil

}

func loadReasons(path string) (*Reasons, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var reasons Reasons
	yaml.Unmarshal(b, &reasons)
	return &reasons, nil
}

func randomElement(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func (s ShoudldIDeploy) isFriday(date time.Time) (bool, string) {
	if date.Weekday() == time.Friday {
		if date.Day() == 13 {
			return true, randomElement(s.reasons.Friday13Th)
		}
		return true, randomElement(s.reasons.FridayAfternoon)
	}
	return false, ""
}

func (s ShoudldIDeploy) isWeekend(date time.Time) (bool, string) {
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return true, randomElement(s.reasons.Weekend)
	}
	return false, ""
}

func (s ShoudldIDeploy) isAfternoon(date time.Time) (bool, string) {
	if date.Hour() >= 16 {
		if date.Weekday() == time.Thursday {
			return true, randomElement(s.reasons.ThursdayAfternoon)
		}
		return true, randomElement(s.reasons.Afternoon)
	}
	return false, ""
}

func (s ShoudldIDeploy) isDayBeforeChristmas(date time.Time) (bool, string) {
	if date.Month() == time.December && date.Day() == 24 {
		return true, randomElement(s.reasons.DayBeforeChristmas)
	}
	return false, ""
}

func (s ShoudldIDeploy) isChristmas(date time.Time) (bool, string) {
	if date.Month() == time.December && (date.Day() == 25 || date.Day() == 26) {
		return true, randomElement(s.reasons.Christmas)
	}
	return false, ""
}

func (s ShoudldIDeploy) isNewYear(date time.Time) (bool, string) {
	if date.Month() == time.December && date.Day() == 31 {
		return true, randomElement(s.reasons.NewYear)
	}

	return false, ""
}

func (s ShoudldIDeploy) isPublicHoliday(date time.Time, country string) (bool, string) {
	if s.holidayService.IsPublicHoliday(date, country) {
		return true, randomElement(s.reasons.ToNotDdeploy) //TODO look for proper reasons
	}
	return false, ""
}

func (s ShoudldIDeploy) ShouldIDeploy(date time.Time, country string) (bool, string) {
	if b, m := s.isFriday(date); b {
		return !b, m
	}
	if b, m := s.isWeekend(date); b {
		return !b, m
	}
	if b, m := s.isAfternoon(date); b {
		return !b, m
	}
	if b, m := s.isDayBeforeChristmas(date); b {
		return !b, m
	}
	if b, m := s.isChristmas(date); b {
		return !b, m
	}
	if b, m := s.isNewYear(date); b {
		return !b, m
	}
	if b, m := s.isPublicHoliday(date, country); b {
		return !b, m
	}

	return true, randomElement(s.reasons.ToDeploy)
}
