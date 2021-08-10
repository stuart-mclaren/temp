package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

type LocationType struct {
	Precipitation Precipitation `xml:" precipitation,omitempty"`
	Temperature   Temperature   `xml:" temperature,omitempty"`
}

type Precipitation struct {
	Unit  string  `xml:"unit,attr"`
	Value float64 `xml:"value,attr"`
}

type Temperature struct {
	Unit  string  `xml:"unit,attr"`
	Value float64 `xml:"value,attr"`
}

type ProductType struct {
	Time []TimeType `xml:" time"`
}

type TimeType struct {
	Location []LocationType `xml:" location"`
	From     time.Time      `xml:"from,attr"`
	To       time.Time      `xml:"to,attr"`
}

type Unitvalue struct {
	Unit  string  `xml:"unit,attr"`
	Value float32 `xml:"value,attr"`
}

type Weatherdata struct {
	Product []ProductType `xml:" product,omitempty"`
}

func findMinAndMaxTemperature(w Weatherdata, hours int,
	now time.Time, zone string) (float64, float64) {
	current := 0
	tz, err := time.LoadLocation(zone)
	if err != nil {
		log.Fatal(`Failed to load location "Local" for timezone`)
	}
	var pmin *float64 = nil
	var pmax *float64 = nil
	for _, t := range w.Product[0].Time {
		// Filter out shorter time difference data
		if t.To.Sub(t.From).Hours() != 0 {
			continue
		}
		// Filter out anything stale (predicting the past)
		if t.To.In(tz).Sub(now).Hours() < 0 {
			continue
		}
		for _, l := range t.Location {
			if pmin == nil {
				pmin = new(float64)
				*pmin = l.Temperature.Value
			} else if l.Temperature.Value < *pmin {
				*pmin = l.Temperature.Value
			}
			if pmax == nil {
				pmax = new(float64)
				*pmax = l.Temperature.Value
			} else if l.Temperature.Value > *pmax {
				*pmax = l.Temperature.Value
			}
			current++
			if current == hours {
				return *pmin, *pmax
			}
		}
	}
	return *pmin, *pmax
}

func display(f io.Writer, w Weatherdata, hours int,
	now time.Time, zone string, min float64, max float64) {
	current := 0
	tz, err := time.LoadLocation(zone)
	if err != nil {
		log.Fatal(`Failed to load location "Local" for timezone`)
	}
	for _, t := range w.Product[0].Time {
		// Filter out shorter time difference data
		if t.To.Sub(t.From).Hours() != 0 {
			continue
		}
		// Filter out anything stale (predicting the past)
		if t.To.In(tz).Sub(now).Hours() < 0 {
			continue
		}
		for _, l := range t.Location {
			temp := int(math.Round(l.Temperature.Value))
			abstemp := int(math.Abs(math.Round(l.Temperature.Value)))
			var offset string
			var stars string
			if temp >= 0 {
				stars = "|" + strings.Repeat("-", temp)
			} else {
				stars = strings.Repeat("-", abstemp) + "|"
			}
			offset = strings.Repeat(" ", int(math.Abs(math.Min(0.0, math.Round(min))))+
				int(math.Min(0.0, float64(temp))))

			fmt.Fprintf(f, "%02d:%02d %3d°C %s%s\n",
				t.From.In(tz).Hour(), t.From.In(tz).Minute(),
				temp,
				offset, stars,
			)
			current++
			if current == hours {
				return
			}
		}
	}
}

func main() {

	w := Weatherdata{}

	hours := flag.Int("hours", 12, "Number of hours to forecast")
	latitude := flag.Float64(
		"latitude", 53.292148,
		"Latitude. Use to specify forecast location",
	)
	longitude := flag.Float64(
		"longitude", -9.007064,
		"Longitude. Use to specify forecast location",
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s: Rain forecast\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Note: Single or double dash can be used "+
			"for parameters, eg -hours/--hours.\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	url := "https://api.met.no/weatherapi/locationforecast/2.0/classic" +
		"?lat=" + fmt.Sprintf("%f", *latitude) +
		";lon=" + fmt.Sprintf("%f", *longitude)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "curl/7.64.0")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = xml.Unmarshal(body, &w)
	if err != nil {
		log.Fatal(err)
	}

	t := time.Now()
	min, max := findMinAndMaxTemperature(w, *hours, t, "Local")
	display(os.Stdout, w, *hours, t, "Local", min, max)
}
