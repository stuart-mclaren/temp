# temp: what will the temperature be in the next few hours?

`temp` is a CLI tool which uses  `yr.no` weather [APIs](http://om.yr.no/verdata/free-weather-data/) to display a temperature forecast for the next several hours.

## Usage

```
$ temp
22:00  10°C      |----------
23:00   9°C      |---------
00:00   7°C      |-------
01:00   3°C      |---
02:00   0°C      |
03:00  -1°C     -|
04:00  -2°C    --|
05:00  -5°C -----|
06:00  -4°C  ----|
07:00  -3°C   ---|
08:00  -1°C     -|
09:00   1°C      |-
```

The above output shows expected temperature per hour over the coming hours.

`temp` assumes a default longitude and latitude. You will usually need to override these to match your location.

```
$ temp --latitude 43.2648 --longitude -18.9297
```

Type `temp --help` for a full list of options.

## Installation

From source

```
$ git clone https://github.com/stuart-mclaren/temp
$ cd temp
$ go build ./...
$ ./temp
```
