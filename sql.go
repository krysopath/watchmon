package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v2"
)

// CreateTableStmt does not need no stinking migrations
const CreateTableStmt string = `
	CREATE TABLE batteryinfo (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	charge_now INTEGER NULL,
	charge_full INTEGER NULL,
	charge_design INTEGER NULL,
	current_now INTEGER NULL,
	voltage_now INTEGER NULL,
	charging INTEGER NULL,
	timestamp INTEGER NULL,
	cycles INTEGER NULL
    );`

//BatteryDataRow holds a measurement
type BatteryDataRow struct {
	Id               int64
	ChargeNow        int64
	ChargeFull       int64
	ChargeFullDesign int64
	CurrentNow       int64
	VoltageNow       int64
	Charging         int64
	Timestamp        int64
	Cycles           int64
}

func (bdw *BatteryDataRow) GetPower() float64 {
	return float64(bdw.VoltageNow / 1000000 * bdw.CurrentNow / 1000000)
}

func (bdw *BatteryDataRow) GetCurrentNow() float64 {
	return float64(bdw.CurrentNow / 1000)
}
func (bdw *BatteryDataRow) GetVoltageNow() float64 {
	return float64(bdw.VoltageNow / 1000)
}

func (bdw *BatteryDataRow) GetCharging() bool {
	if bdw.Charging == 0 {
		return false
	} else {
		return true
	}
}

func (bdw *BatteryDataRow) GetChargeNow() float64 {
	return float64(bdw.ChargeNow / 1000)
}

func (bdw *BatteryDataRow) GetChargeFull() float64 {
	return float64(bdw.ChargeFull / 1000)
}

func (bdw *BatteryDataRow) GetChargeFullDesign() float64 {
	return float64(bdw.ChargeFullDesign / 1000)
}

func (bdw *BatteryDataRow) GetCapacityPermille() int64 {
	return 1000.0 * bdw.ChargeNow / bdw.ChargeFull
}

func (bdw *BatteryDataRow) GetCapacityDegradation() int64 {
	return -1000.0*bdw.ChargeFull/bdw.ChargeFullDesign + 1000
}

func (bdw *BatteryDataRow) Compute() *BatteryDataComputed {
	return &BatteryDataComputed{
		bdw.Id,
		bdw.Timestamp,
		bdw.GetCharging(),
		bdw.GetPower(),
		bdw.GetCurrentNow(),
		bdw.GetVoltageNow(),
		bdw.GetChargeNow(),
		bdw.GetChargeFull(),
		bdw.GetChargeFullDesign(),
		bdw.GetCapacityPermille(),
		bdw.GetCapacityDegradation(),
		bdw.Cycles,
	}
}

type BatteryDataComputed struct {
	Id                  int64   `json:"id" yaml:"id"`
	Timestamp           int64   `json:"timestamp" yaml:"timestamp"`
	Charging            bool    `json:"charging" yaml:"charging"`
	Power               float64 `json:"power_W" yaml:"power_W"`
	CurrentNow          float64 `json:"current_now_mA" yaml:"current_now_mA"`
	VoltageNow          float64 `json:"voltage_now_mV" yaml:"voltage_now_mV"`
	ChargeNow           float64 `json:"charge_now_mAh" yaml:"charge_now_mAh"`
	ChargeFull          float64 `json:"charge_full_mAh" yaml:"charge_full_mAh"`
	ChargeFullDesign    float64 `json:"charge_full_design_mAh" yaml:"charge_full_design_mAh"`
	CapacityPermille    int64   `json:"capacity_permille" yaml:"capacity_permille"`
	CapacityDegradation int64   `json:"capacity_degration_permille" yaml:"capacity_degration_permille"`
	Cycles              int64   `json:"cycles" yaml:"cycles"`
}

func (bdc *BatteryDataComputed) String() string {
	fString := `        Timestamp: %10d
         Charging: %10t
            Power: %10.1f W
       CurrentNow: %10.1f mA
       VoltageNow: %10.1f mV
        ChargeNow: %10.1f mAh
       ChargeFull: %10.1f mAh
 ChargeFullDesign: %10.1f mAh
 CapacityPermille: %10d ‰
CapacityDegration: %10d ‰
           Cycles: %10d ‰
`
	return fmt.Sprintf(
		fString,
		bdc.Timestamp,
		bdc.Charging,
		bdc.Power,
		bdc.CurrentNow,
		bdc.VoltageNow,
		bdc.ChargeNow,
		bdc.ChargeFull,
		bdc.ChargeFullDesign,
		bdc.CapacityPermille,
		bdc.CapacityDegradation,
		bdc.Cycles,
	)
}

func (bdc *BatteryDataComputed) Json() string {
	asBytes, err := json.Marshal(bdc)
	checkErr(err)
	return string(asBytes)
}

func (bdc *BatteryDataComputed) Yaml() string {
	asBytes, err := yaml.Marshal(bdc)
	checkErr(err)
	return string(asBytes)
}

// CreateDatabaseAndTables or bust!
func CreateDatabaseAndTables(db *sql.DB) {
	_, err := db.Exec(CreateTableStmt)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
