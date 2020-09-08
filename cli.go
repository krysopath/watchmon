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
)

func ParseFlags() *Cli {
	batteryDevice := flag.String(
		"bat", "BAT0", "specify battery device from '/sys/class/power_supply/BAT*'")

	return &Cli{
		BatteryDevice: batteryDevice,
		SqliteFile:    &DBFileDefault,
	}

}

type Cli struct {
	BatteryDevice  *string
	SqliteFile     *string
	CreateDBToggle *bool
	DumpRowsToggle *bool
	OutputFormat   *string
	DB             *sql.DB
	Last           *BatteryDataRow
}

func (c *Cli) String() string {
	cfgBytes, err := yaml.Marshal(c)
	checkErr(err)
	return string(cfgBytes)
}

func (cli *Cli) Init() {
	db, err := sql.Open("sqlite3", *cli.SqliteFile)
	checkErr(err)
	cli.DB = db

	//if *cli.CreateDBToggle {
	//	CreateDatabaseAndTables(cli.DB)
	//}
}

func (cli *Cli) DumpRows() {
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

	fmt.Println("id|charge_now|charge_full|charge_design|current_now|voltage_now|charging|timestamp")
	for rows.Next() {
		err = rows.Scan(
			&id,
			&chargeNow,
			&chargeFull,
			&chargeDesign,
			&currentNow,
			&voltageNow,
			&charging,
			&timestamp)

		checkErr(err)
		fmt.Printf(
			"%d|%d|%d|%d|%d|%d|%d|%d\n",
			id,
			chargeNow,
			chargeFull,
			chargeDesign,
			currentNow,
			voltageNow,
			charging,
			timestamp)
	}

}

func (cli *Cli) FormatRow(bdw *BatteryDataRow) string {
	computed := bdw.Compute()
	switch *cli.OutputFormat {
	case "plain":
		return computed.String()
	case "yaml":
		return computed.Yaml()
	case "json":
		return computed.Json()
	default:
		return computed.String()
	}
}

func (cli *Cli) Measure() string {
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

	var batInfo *BatteryDataRow = CreateBatteryData(cli.BatteryDevice)

	_, errSQL := stmt.Exec(
		batInfo.ChargeNow,
		batInfo.ChargeFull,
		batInfo.ChargeFullDesign,
		batInfo.CurrentNow,
		batInfo.VoltageNow,
		batInfo.Charging,
		batInfo.Timestamp,
	)
	checkErr(errSQL)
	return cli.FormatRow(batInfo)
}

func (cli *Cli) Completions() string {
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
	measureCmd := flag.NewFlagSet("measure", flag.ExitOnError)
	batteryDevice := measureCmd.String(
		"bat",
		"BAT0",
		"specify battery device from '/sys/class/power_supply/BAT*'",
	)
	outputFormat := measureCmd.String("output", "plain", "choose from plain,json,yaml")
	dbCmd := flag.NewFlagSet("db", flag.ExitOnError)
	completionCmd := flag.NewFlagSet("completion", flag.ExitOnError)

	cli := &Cli{
		BatteryDevice: batteryDevice,
		SqliteFile:    &DBFileDefault,
		OutputFormat:  outputFormat,
	}

	cli.Init()
	defer cli.DB.Close()

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "measure":
			measureCmd.Parse(os.Args[2:])
			fmt.Fprintf(os.Stdout, "%+v", cli.Measure())
		case "dump":
			dbCmd.Parse(os.Args)
			fmt.Println(dbCmd.Args())
			cli.DumpRows()
		case "completions":
			completionCmd.Parse(os.Args)
			fmt.Fprintf(os.Stdout,
				"%s\n",
				cli.Completions())
		default:
			usage()
			os.Exit(1)
		}
	} else {
		usage()
	}
}
