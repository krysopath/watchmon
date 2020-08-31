package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

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

func (bdw *BatteryDataRow) String() string {
	fString := `
        Timestamp: %10d
            Power: %10.1f W
       CurrentNow: %10.1f mA
       VoltageNow: %10.1f mV
         Charging: %10b
        ChargeNow: %10.1f mAh
       ChargeFull: %10.1f mAh
 ChargeFullDesign: %10.1f mAh
 CapacityPermille: %10d ‰
CapacityDegration: %10d ‰
`
	return fmt.Sprintf(
		fString,
		bdw.Timestamp,
		float64(bdw.VoltageNow/1000000*bdw.CurrentNow/1000000),
		float64(bdw.CurrentNow/1000),
		float64(bdw.VoltageNow/1000),
		bdw.Charging,
		float64(bdw.ChargeNow/1000),
		float64(bdw.ChargeFull/1000),
		float64(bdw.ChargeFullDesign/1000),
		1000.0*bdw.ChargeNow/bdw.ChargeFull,
		-1000.0*bdw.ChargeFull/bdw.ChargeFullDesign+1000,
	)
}

// CreateDatabaseAndTables
func CreateDatabaseAndTables(db *sql.DB) {
	_, err := db.Exec(CreateTableStmt)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
