package main

import (
	"github.com/abaril/GoLights/api"
	"github.com/cheekybits/is"
	"testing"
	"time"
)

func TestIsCurrentlyDarkOutside(t *testing.T) {

	is := is.New(t)
	db := api.NewMemDB()

	dark := isCurrentlyDarkOutside(db, time.Date(2015, time.August, 4, 12, 0, 0, 0, time.Local))
	is.Equal(dark, true)

	f := WeatherForecast{
		SunriseTime: Timestamp{Time: time.Date(2015, time.August, 4, 8, 0, 0, 0, time.Local)},
		SunsetTime:  Timestamp{Time: time.Date(2015, time.August, 4, 20, 0, 0, 0, time.Local)},
	}
	db.Set("Weather", f)
	dark = isCurrentlyDarkOutside(db, time.Date(2015, time.August, 4, 12, 0, 0, 0, time.Local))
	is.Equal(dark, false)

	dark = isCurrentlyDarkOutside(db, time.Date(2015, time.August, 4, 7, 58, 0, 0, time.Local))
	is.Equal(dark, true)

	dark = isCurrentlyDarkOutside(db, time.Date(2015, time.August, 4, 20, 02, 0, 0, time.Local))
	is.Equal(dark, true)

}
