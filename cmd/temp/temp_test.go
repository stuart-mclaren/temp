package main

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"testing"
	"time"
)

func TestDisplay(t *testing.T) {

	body, err := ioutil.ReadFile("../../test/fixture.xml")
	if err != nil {
		t.Errorf("Could not read file")
	}
	w := Weatherdata{}
	err = xml.Unmarshal(body, &w)
	if err != nil {
		t.Errorf("Could not marshall xml")
	}
	buf := new(bytes.Buffer)
	location, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		panic(err)
	}
	now := time.Date(
		2019, 5, 9, 12, 34, 58, 651387237, location)
	min, max := findMinAndMaxTemperature(w, 12, now, "Europe/Dublin")
	display(buf, w, 12, now, "Europe/Dublin", min, max)
	expected := `13:00  20°C                      |--------------------
14:00 -21°C ---------------------|
15:00  -6°C                ------|
16:00  41°C                      |-----------------------------------------
17:00   0°C                      |
18:00  11°C                      |-----------
19:00  10°C                      |----------
20:00  10°C                      |----------
21:00   9°C                      |---------
22:00   8°C                      |--------
23:00   8°C                      |--------
00:00   7°C                      |-------
`
	if buf.String() != expected {
		t.Errorf("Output mismatch. Expected %s, Actual %s",
			expected, buf.String())
	}
}

func TestDisplaySkipFirstHour(t *testing.T) {

	body, err := ioutil.ReadFile("../../test/fixture.xml")
	if err != nil {
		t.Errorf("Could not read file")
	}
	w := Weatherdata{}
	err = xml.Unmarshal(body, &w)
	if err != nil {
		t.Errorf("Could not marshall xml")
	}
	buf := new(bytes.Buffer)
	location, err := time.LoadLocation("Europe/Dublin")
	if err != nil {
		panic(err)
	}
	now := time.Date(
		2019, 5, 9, 13, 34, 58, 651387237, location)
	min, max := findMinAndMaxTemperature(w, 18, now, "Europe/Dublin")
	display(buf, w, 18, now, "Europe/Dublin", min, max)
	expected := `14:00 -21°C ---------------------|
15:00  -6°C                ------|
16:00  41°C                      |-----------------------------------------
17:00   0°C                      |
18:00  11°C                      |-----------
19:00  10°C                      |----------
20:00  10°C                      |----------
21:00   9°C                      |---------
22:00   8°C                      |--------
23:00   8°C                      |--------
00:00   7°C                      |-------
01:00   6°C                      |------
02:00   6°C                      |------
03:00   6°C                      |------
04:00   6°C                      |------
05:00   6°C                      |------
06:00   6°C                      |------
07:00   6°C                      |------
`
	if buf.String() != expected {
		t.Errorf("Output mismatch. Expected %s, Actual %s",
			expected, buf.String())
	}
}

func TestDisplayNonHourTimezone(t *testing.T) {

	body, err := ioutil.ReadFile("../../test/fixture.xml")
	if err != nil {
		t.Errorf("Could not read file")
	}
	w := Weatherdata{}
	err = xml.Unmarshal(body, &w)
	if err != nil {
		t.Errorf("Could not marshall xml")
	}
	buf := new(bytes.Buffer)
	location, err := time.LoadLocation("Asia/Rangoon")
	if err != nil {
		panic(err)
	}
	now := time.Date(
		2019, 5, 9, 13, 34, 58, 651387237, location)
	min, max := findMinAndMaxTemperature(w, 18, now, "Asia/Rangoon")
	display(buf, w, 18, now, "Asia/Rangoon", min, max)
	expected := `18:30  20°C                      |--------------------
19:30 -21°C ---------------------|
20:30  -6°C                ------|
21:30  41°C                      |-----------------------------------------
22:30   0°C                      |
23:30  11°C                      |-----------
00:30  10°C                      |----------
01:30  10°C                      |----------
02:30   9°C                      |---------
03:30   8°C                      |--------
04:30   8°C                      |--------
05:30   7°C                      |-------
06:30   6°C                      |------
07:30   6°C                      |------
08:30   6°C                      |------
09:30   6°C                      |------
10:30   6°C                      |------
11:30   6°C                      |------
`
	if buf.String() != expected {
		t.Errorf("Output mismatch. Expected %s, Actual %s",
			expected, buf.String())
	}
}
