package main

import (
	"database/sql"

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

func CreateDatabaseAndTables(db *sql.DB) {
	_, err := db.Exec(CreateTableStmt)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
