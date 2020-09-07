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

func CreateBatteryData(batteryDevice *string) *BatteryDataRow {
	var batteryPath string = "/sys/class/power_supply/%s/%s"
	return &BatteryDataRow{
		ChargeNow: ReadBatteryValueAsInt(
			batteryPath, *batteryDevice, "charge_now"),
		ChargeFull: ReadBatteryValueAsInt(
			batteryPath, *batteryDevice, "charge_full"),
		ChargeFullDesign: ReadBatteryValueAsInt(
			batteryPath, *batteryDevice, "charge_full_design"),
		CurrentNow: ReadBatteryValueAsInt(
			batteryPath, *batteryDevice, "current_now"),
		VoltageNow: ReadBatteryValueAsInt(
			batteryPath, *batteryDevice, "voltage_now"),
		Charging: IsBatteryCharging(
			batteryPath, *batteryDevice, "status"),
		Timestamp: time.Now().Unix(),
	}
}
