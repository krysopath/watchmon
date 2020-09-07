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
	timestamp INTEGER NULL
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

func (bdw *BatteryDataRow) GetCharging() string {
	if bdw.Charging == 0 {
		return "no"
	} else {
		return "yes"
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

func (bdw *BatteryDataRow) String() string {
	fString := `        Timestamp: %10d
            Power: %10.1f W
       CurrentNow: %10.1f mA
       VoltageNow: %10.1f mV
         Charging: %10s
        ChargeNow: %10.1f mAh
       ChargeFull: %10.1f mAh
 ChargeFullDesign: %10.1f mAh
 CapacityPermille: %10d ‰
CapacityDegration: %10d ‰
`
	return fmt.Sprintf(
		fString,
		bdw.Timestamp,
		bdw.GetPower(),
		bdw.GetCurrentNow(),
		bdw.GetVoltageNow(),
		bdw.GetCharging(),
		bdw.GetChargeNow(),
		bdw.GetChargeFull(),
		bdw.GetChargeFullDesign(),
		bdw.GetCapacityPermille(),
		bdw.GetCapacityDegradation(),
	)
}

func (bdw *BatteryDataRow) Json() string {
	cfgBytes, err := json.Marshal(bdw)
	checkErr(err)
	return string(cfgBytes)
}

func (bdw *BatteryDataRow) Yaml() string {
	cfgBytes, err := yaml.Marshal(bdw)
	checkErr(err)
	return string(cfgBytes)
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
