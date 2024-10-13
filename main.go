package main

import (
	"go-shouldideploy/pkg/api"
	"net/http"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	s, err := api.NewShouldIDeploy("reasons.yaml")
	if err != nil {
		panic(err)
	}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			langs := strings.Split(c.Request().Header.Get("Accept-Language"), ",")
			if len(langs) > 0 {
				lang := strings.Split(langs[0], "-")
				if len(lang) > 1 {
					c.Set("country", lang[1])
				} else {
					c.Set("country", "DE")
				}

			} else {
				c.Set("country", "DE")
			}
			return next(c)

		}
	})
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
			date, err = time.Parse(time.DateOnly, dateParam)
			if err != nil {
				return c.String(http.StatusBadRequest, err.Error())
			}
		}
		b, m := s.ShouldIDeploy(date.In(tz), c.Get("country").(string))

		return c.JSON(http.StatusOK, map[string]interface{}{"shouldIDeploy": b, "message": m})
	})
	e.Static("/", "static")
	e.Logger.Fatal(e.Start(":3000"))
}
