package main

import (
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"gopkg.in/yaml.v2"
)

func getUser() *user.User {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr
}

var (
	User          = getUser()
	DBFileDefault = fmt.Sprintf(
		"%s/watchmon.sqlite",
		User.HomeDir)
	bashCompletion string
	gitTag         string
	gitRef         string
	outputFormat   string
	Cli            *CLI

	measureCmd    = flag.NewFlagSet("measure", flag.ExitOnError)
	batteryDevice = measureCmd.String(
		"bat",
		"BAT0",
		"specify battery device from '/sys/class/power_supply/BAT*'",
	)
	intendedFormat = measureCmd.String("output", "plain", "choose from plain,json,yaml")
	noStore        = measureCmd.Bool("no-store", false, "choose to not store a measurement")

	dbCmd    = flag.NewFlagSet("db", flag.ExitOnError)
	dbCreate = flag.Bool("create", false, "wether to create the sqlite file")

	completionCmd = flag.NewFlagSet("completion", flag.ExitOnError)
)

type CLI struct {
	BatteryDevice  *string
	SqliteFile     *string
	CreateDBToggle *bool
	DumpRowsToggle *bool
	NoStoreToggle  *bool
	OutputFormat   *string
	DB             *sql.DB
	Last           *BatteryDataRow
}

func (c *CLI) String() string {
	cfgBytes, err := yaml.Marshal(c)
	checkErr(err)
	return string(cfgBytes)
}

func (cli *CLI) Init() {
	db, err := sql.Open("sqlite3", *cli.SqliteFile)
	checkErr(err)
	cli.DB = db

	if *cli.CreateDBToggle {
		CreateDatabaseAndTables(cli.DB)
	}
}

func (cli *CLI) DumpRows() {
	rows, err := cli.DB.Query("SELECT * FROM batteryinfo")
	checkErr(err)

	var id int
	var chargeNow int
	var chargeFull int
	var chargeDesign int
	var currentNow int
	var voltageNow int
	var charging int
	var timestamp int
	var cycleCount int

	fmt.Println("id|charge_now|charge_full|charge_design|current_now|voltage_now|charging|timestamp|cycle_count")
	for rows.Next() {
		err = rows.Scan(
			&id,
			&chargeNow,
			&chargeFull,
			&chargeDesign,
			&currentNow,
			&voltageNow,
			&charging,
			&timestamp,
			&cycleCount,
		)

		checkErr(err)
		fmt.Printf(
			"%d|%d|%d|%d|%d|%d|%d|%d|%d\n",
			id,
			chargeNow,
			chargeFull,
			chargeDesign,
			currentNow,
			voltageNow,
			charging,
			timestamp,
			cycleCount,
		)
	}

}

func (cli *CLI) Measure() *BatteryDataComputed {

	var batInfo *BatteryDataRow = CreateBatteryData(cli.BatteryDevice)

	if !*cli.NoStoreToggle {

		stmt, err := cli.DB.Prepare(`
		INSERT INTO batteryinfo(
			charge_now, 
			charge_full, 
			charge_design, 
			current_now, 
			voltage_now, 
			charging, 
			timestamp,
			cycle_count) 
		values(?,?,?,?,?,?,?,?)`,
		)
		checkErr(err)

		_, errSQL := stmt.Exec(
			batInfo.ChargeNow,
			batInfo.ChargeFull,
			batInfo.ChargeFullDesign,
			batInfo.CurrentNow,
			batInfo.VoltageNow,
			batInfo.Charging,
			batInfo.Timestamp,
			batInfo.Cycles,
		)
		checkErr(errSQL)
	}
	return batInfo.Compute()
}

func (cli *CLI) LastMeasure() *BatteryDataComputed {
	var id int64
	var chargeNow int64
	var chargeFull int64
	var chargeDesign int64
	var currentNow int64
	var voltageNow int64
	var charging int64
	var timestamp int64
	var cycleCount int64

	err := cli.DB.QueryRow(
		"SELECT * FROM batteryinfo ORDER BY id DESC LIMIT 1",
	).Scan(
		&id,
		&chargeNow,
		&chargeFull,
		&chargeDesign,
		&currentNow,
		&voltageNow,
		&charging,
		&timestamp,
		&cycleCount,
	)
	checkErr(err)
	bdw := BatteryDataRow{
		Id:               id,
		ChargeNow:        chargeNow,
		ChargeFull:       chargeFull,
		ChargeFullDesign: chargeDesign,
		CurrentNow:       currentNow,
		VoltageNow:       voltageNow,
		Charging:         charging,
		Timestamp:        timestamp,
		Cycles:           cycleCount,
	}
	return bdw.Compute()
}

func (cli *CLI) Completions() string {
	data, err := base64.StdEncoding.DecodeString(bashCompletion)
	if err != nil {
		panic("error: the shell completions script could not be decoded")
	}
	return string(data)

}

func usage() {
	fmt.Fprintf(
		os.Stderr,
		`watchmon %s-%s
		supports measure, dump, completions sub commands`,
		gitTag, gitRef,
	)
}

func main() {
	flag.Parse()
	Cli = &CLI{
		BatteryDevice:  batteryDevice,
		SqliteFile:     &DBFileDefault,
		NoStoreToggle:  noStore,
		CreateDBToggle: dbCreate,
		OutputFormat:   intendedFormat,
	}
	Cli.Init()

	defer Cli.DB.Close()

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "measure":
			measureCmd.Parse(os.Args[2:])
			Cli.Measure()
			fmt.Fprintf(os.Stdout, "%+v", Cli.LastMeasure())
		case "dump":
			dbCmd.Parse(os.Args)
			fmt.Println(dbCmd.Args())
			Cli.DumpRows()
		case "completions":
			completionCmd.Parse(os.Args)
			fmt.Fprintf(os.Stdout,
				"%s\n",
				Cli.Completions())
		default:
			usage()
			os.Exit(1)
		}
	} else {
		usage()
	}
}
