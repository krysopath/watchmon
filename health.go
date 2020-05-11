package main

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func ReadBatteryValue(path string) int64 {
	asBytes, err := ioutil.ReadFile(path)
	checkErr(err)
	asString := strings.TrimSpace(string(asBytes))
	asInt, err := strconv.ParseInt(asString, 10, 64)
	checkErr(err)
	return asInt
}

func Measure() *BatteryDataRow {
	return &BatteryDataRow{
		ChargeNow:        ReadBatteryValue("/sys/class/power_supply/BAT0/charge_now"),
		ChargeFull:       ReadBatteryValue("/sys/class/power_supply/BAT0/charge_full"),
		ChargeFullDesign: ReadBatteryValue("/sys/class/power_supply/BAT0/charge_full_design"),
		CurrentNow:       ReadBatteryValue("/sys/class/power_supply/BAT0/current_now"),
		VoltageNow:       ReadBatteryValue("/sys/class/power_supply/BAT0/voltage_now"),
		Timestamp:        time.Now().Unix(),
	}
}
