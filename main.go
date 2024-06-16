package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/*
Done
create sqlite on startup if not already existing
check if db already created
create datamodel for dayentry
	- date
	- start time
	- end time
add day entry

Todo
Move everything db related to different file
abstract away db operations

edit day entry
delete day entry
create datamodel for timeentry
	- requires day entry
	- time hh:mm
	- tic key
	- comment
	- story number optional
add timeentry
track time
edit timeentry
delete timeentry

visual presentation when editing
	- navigate with arrowkeys/hjkl
	- try something like the add dialog
*/

const initTable string = `
  CREATE TABLE IF NOT EXISTS dayEntry (
  id INTEGER NOT NULL PRIMARY KEY,
  date DATETIME NOT NULL,
  start string,
  end string
  );
  CREATE TABLE IF NOT EXISTS timeEntry (
  id INTEGER NOT NULL PRIMARY KEY,
  dayEntryId INTEGER,
  duration INTEGER,
  key string,
  story_number string,
  comment string
  );
  `

const enterDatePrompt string = "Please enter a date in yyyy-mm-dd format or press enter for current date \n> "
const enterStartTimePrompt string = "Please enter:\n - a start time in hh:mm format or press\n - none\n - press enter to use current time\n> "
const enterEndTimePrompt string = "Please enter:\n - a end time in hh:mm format or press\n - none\n - press enter to use current time\n> "

type dayEntry struct {
	date  time.Time
	start string
	end   string
}

func main() {
	// DB Conection code
	db, err := sql.Open("sqlite3", "chronos.db")
	exitOnErr(err, db)
	_, err = db.Exec(initTable)
	exitOnErr(err, db)

	// FLAG CODE
	addDayFlag := flag.Bool("addDay", false, "Add a day entry to chronos")
	flag.Parse()
	if *addDayFlag {
		// Read first argument
		buffer := bufio.NewReader(os.Stdin)
		fmt.Print(enterDatePrompt)
		dateEntry, err := buffer.ReadBytes('\n')
		exitOnErr(err, db)
		dateResult, err := parseDate(dateEntry)
		exitOnErr(err, db)

		// Read second argument
		buffer.Reset(os.Stdin)
		fmt.Print(enterStartTimePrompt)
		startTimeEntry, err := buffer.ReadBytes('\n')
		exitOnErr(err, db)
		startTimeResult, err := parseTime(startTimeEntry)
		exitOnErr(err, db)

		// Read third argument
		buffer.Reset(os.Stdin)
		fmt.Print(enterEndTimePrompt)
		endTimeEntry, err := buffer.ReadBytes('\n')
		exitOnErr(err, db)
		endTimeResult, err := parseTime(endTimeEntry)
		exitOnErr(err, db)
		d := dayEntry{
			dateResult,
			startTimeResult,
			endTimeResult,
		}
		res, err := db.Exec("INSERT INTO dayEntry VALUES(null,?,?,?)", d.date, d.start, d.end)
		exitOnErr(err, db)
		fmt.Println(res)
	}

	db.Close()
}

func parseDate(date []byte) (resDate time.Time, err error) {
	var result time.Time
	var resultError error
	var dateRegex string = "^\\d{4}-(\\d{2}|\\d{1})-(\\d{2}|\\d{1})"
	matchDateFormat, err := regexp.Match(dateRegex, date)
	if err != nil {
		resultError = err
	}

	if matchDateFormat {
		r, err := regexp.Compile("-")
		if err != nil {
			resultError = err
		}

		datesGroup := r.Split(string(date), 3)
		if err != nil {
			resultError = err
		}

		result, err = constructDate(datesGroup[0], datesGroup[1], strings.TrimSuffix(datesGroup[2], "\n"))
		if err != nil {
			resultError = err
		}
		return result, resultError
	}

	trimmedDate := strings.TrimSuffix(string(date), "\n")
	if trimmedDate == "" {
		return time.Date(
			time.Now().Local().Year(),
			time.Now().Local().Month(),
			time.Now().Local().Day(),
			0, 0, 0, 0, time.Local), resultError
	}
	resultError = errors.New("provided date did not match the requested format yyyy-mm-dd")
	return result, resultError
}

func constructDate(year string, month string, day string) (resDate time.Time, err error) {
	var result time.Time
	var resultError error

	years, err := strconv.Atoi(year)
	if err != nil {
		resultError = err
	}
	months, err := strconv.Atoi(month)
	if err != nil {
		resultError = err
	}
	days, err := strconv.Atoi(day)
	if err != nil {
		resultError = err
	}
	result = time.Date(years, time.Month(months), days, 0, 0, 0, 0, time.Local)
	return result, resultError
}

func parseTime(timeString []byte) (resTime string, err error) {
	var result string
	var resultError error
	var timeRegex string = "^(\\d{2}|\\d{1}):(\\d{2}|\\d{1})"
	matchTimeRegex, err := regexp.Match(timeRegex, timeString)
	if err != nil {
		resultError = err
	}
	trimmedTimeString := strings.TrimSuffix(string(timeString), "\n")

	if matchTimeRegex {
		result = trimmedTimeString
		// handle format like 4:35 or 15:3 and return error
		return result, resultError
	}

	if trimmedTimeString == "none" {
		return "", resultError
	}
	if trimmedTimeString == "" {
		hours := time.Now().Local().Hour()
		minutes := time.Now().Local().Minute()
		result = fmt.Sprint(hours) + ":" + fmt.Sprint(minutes)
		return result, resultError
	}

	return result, resultError
}

func exitOnErr(err error, db *sql.DB) {
	if err != nil {
		fmt.Println(err)
		db.Close()
		os.Exit(1)
	}
}
