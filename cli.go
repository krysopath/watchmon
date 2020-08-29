package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os/user"
)

func getUser() *user.User {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr
}

var User = getUser()

var DBFileDefault = fmt.Sprintf(
	"%s/watchmon.sqlite",
	User.HomeDir)

func ParseFlags() *Cli {
	batteryDevice := flag.String(
		"bat", "BAT0", "specify battery device from '/sys/class/power_supply/BAT*'")
	sqliteFile := flag.String(
		"sqlite", DBFileDefault, "saving data to sqlite")
	dbCreateToggle := flag.Bool(
		"dbcreate", false, "create the table?")
	dumpRowsToggle := flag.Bool(
		"dump", false, "dump the table?")
	flag.Parse()

	return &Cli{
		BatteryDevice:  batteryDevice,
		SqliteFile:     sqliteFile,
		CreateDBToggle: dbCreateToggle,
		DumpRowsToggle: dumpRowsToggle,
	}

}

type Cli struct {
	BatteryDevice  *string
	SqliteFile     *string
	CreateDBToggle *bool
	DumpRowsToggle *bool
	DB             *sql.DB
}

func (cli *Cli) Init() {
	db, err := sql.Open("sqlite3", *cli.SqliteFile)
	checkErr(err)
	cli.DB = db

	if *cli.CreateDBToggle {
		CreateDatabaseAndTables(cli.DB)
	}
}

func (cli *Cli) DumpRows() {
	rows, err := cli.DB.Query("SELECT * FROM batteryinfo")
	checkErr(err)

	var id int
	var charge_now int
	var charge_full int
	var charge_design int
	var current_now int
	var voltage_now int
	var charging int
	var timestamp int

	fmt.Println("id|charge_now|charge_full|charge_design|current_now|voltage_now|charging|timestamp")
	for rows.Next() {
		err = rows.Scan(
			&id,
			&charge_now,
			&charge_full,
			&charge_design,
			&current_now,
			&voltage_now,
			&charging,
			&timestamp)

		checkErr(err)
		fmt.Printf(
			"%d|%d|%d|%d|%d|%d|%d|%d\n",
			id,
			charge_now,
			charge_full,
			charge_design,
			current_now,
			voltage_now,
			charging,
			timestamp)
	}

}

func (cli *Cli) Do() {
	stmt, err := cli.DB.Prepare(`
		INSERT INTO batteryinfo(
			charge_now, 
			charge_full, 
			charge_design, 
			current_now, 
			voltage_now, 
			charging, 
			timestamp) 
		values(?,?,?,?,?,?,?)`,
	)
	checkErr(err)
	var batInfo *BatteryDataRow
	batInfo = Measure(cli.BatteryDevice)

	fmt.Printf("%+v", batInfo)

	res, err := stmt.Exec(
		batInfo.ChargeNow,
		batInfo.ChargeFull,
		batInfo.ChargeFullDesign,
		batInfo.CurrentNow,
		batInfo.VoltageNow,
		batInfo.Charging,
		batInfo.Timestamp,
	)
	checkErr(err)
	fmt.Printf("%+v", res)
}

func main() {
	cli := ParseFlags()
	cli.Init()
	defer cli.DB.Close()

	if *cli.DumpRowsToggle {
		cli.DumpRows()
	} else {
		cli.Do()
	}

}
