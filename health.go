package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

func IsBatteryCharging(path, battery, property string) int64 {
	asBytes, err := ioutil.ReadFile(
		fmt.Sprintf(path, battery, property),
	)
	checkErr(err)
	asString := strings.TrimSpace(string(asBytes))
	if asString == "Discharging" {
		return 0
	} else {
		return 1
	}
}

func ReadBatteryValueAsInt(path, battery, property string) int64 {
	asBytes, err := ioutil.ReadFile(
		fmt.Sprintf(path, battery, property),
	)
	checkErr(err)
	asString := strings.TrimSpace(string(asBytes))
	asInt, err := strconv.ParseInt(asString, 10, 64)
	checkErr(err)
	return asInt
}

func Measure() *BatteryDataRow {
	var batteryPath string = "/sys/class/power_supply/%s/%s"
	return &BatteryDataRow{
		ChargeNow: ReadBatteryValueAsInt(
			batteryPath, "BAT0", "charge_now"),
		ChargeFull: ReadBatteryValueAsInt(
			batteryPath, "BAT0", "charge_full"),
		ChargeFullDesign: ReadBatteryValueAsInt(
			batteryPath, "BAT0", "charge_full_design"),
		CurrentNow: ReadBatteryValueAsInt(
			batteryPath, "BAT0", "current_now"),
		VoltageNow: ReadBatteryValueAsInt(
			batteryPath, "BAT0", "voltage_now"),
		Charging: IsBatteryCharging(
			batteryPath, "BAT0", "status"),
		Timestamp: time.Now().Unix(),
	}
}
