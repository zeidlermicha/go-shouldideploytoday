package main

import (
	"net/http"
	"os"
	"time"
	_ "time/tzdata"

	"github.com/labstack/echo/v4"
	"golang.org/x/exp/rand"
	"gopkg.in/yaml.v3"
)

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

var reasons Reasons

func randomElement(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func isFriday(date time.Time) (bool, string) {
	if date.Weekday() == time.Friday {
		if date.Day() == 13 {
			return true, randomElement(reasons.Friday13Th)
		}
		return true, randomElement(reasons.FridayAfternoon)
	}
	return false, ""
}

func isWeekend(date time.Time) (bool, string) {
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		return true, randomElement(reasons.Weekend)
	}
	return false, ""
}

func isAfternoon(date time.Time) (bool, string) {
	if date.Hour() >= 16 {
		return true, randomElement(reasons.Afternoon)
	}
	return false, ""
}

func shouldIDeploy(date time.Time) (bool, string) {
	if b, m := isFriday(date); b {
		return !b, m
	}
	if b, m := isWeekend(date); b {
		return !b, m
	}
	if b, m := isAfternoon(date); b {
		return !b, m
	}

	return true, randomElement(reasons.ToDeploy)
}

func LoadConfig(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	yaml.Unmarshal(b, &reasons)
	return nil
}

func main() {
	err := LoadConfig("reasons.yaml")
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.GET("/api", func(c echo.Context) error {
		timeZone := c.QueryParam("tz")
		dateParam := c.QueryParam("date")
		if timeZone == "" {
			timeZone = "UTC"
		}
		tz, err := time.LoadLocation(timeZone)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		date := time.Now()
		if dateParam != "" {
			date, err = time.Parse("", dateParam)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
		}
		b, m := shouldIDeploy(date.In(tz))

		return c.JSON(http.StatusOK, map[string]interface{}{"shouldIDeploy": b, "message": m})
	})
	e.Static("/", "static")
	e.Logger.Fatal(e.Start(":3000"))
}
