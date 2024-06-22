package repository

import (
	"chronos/model"
	"database/sql"
	"log"
	"os"
	"time"
)

const initTable string = `
  CREATE TABLE IF NOT EXISTS dayEntry (
  id INTEGER NOT NULL PRIMARY KEY,
  date DATETIME NOT NULL UNIQUE,
  start DATETIME NOT NULL,
  end DATETIME NOT NULL
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

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "chronos.db")
	exitOnErr(err, db)
	_, err = db.Exec(initTable)
	exitOnErr(err, db)
	return db
}

func GetDayEntryIdByDate(db *sql.DB, date time.Time) (*int, error) {
	var resultId int
	row := db.QueryRow("SELECT id from dayEntry where date=?", date)
	err := row.Scan(&resultId)
	if err != nil {
		return nil, err
	}
	return &resultId, nil
}

func GetDayEntryByDate(db *sql.DB, date time.Time) (*int, *model.DayEntry, error) {
	var resultId int
	var result model.DayEntry
	row := db.QueryRow("SELECT id, date, start, end from dayEntry where date=?", date)
	err := row.Scan(&resultId, &result.Date, &result.Start, &result.End)
	if err != nil {
		return nil, nil, err
	}
	return &resultId, &result, nil
}

func GetAllDayEntries(db *sql.DB) (*[]model.DayEntry, error) {
	rows, err := db.Query("SELECT date, start, end from dayEntry ORDER BY id DESC LIMIT 100")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.DayEntry{}
	for rows.Next() {
		i := model.DayEntry{}
		err := rows.Scan(&i.Date, &i.Start, &i.End)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return &result, nil
}

func InsertDayEntry(db *sql.DB, dayEntry *model.DayEntry) {
	_, err := db.Exec("INSERT INTO dayEntry VALUES(null,?,?,?)", dayEntry.Date, dayEntry.Start, dayEntry.End)
	exitOnErr(err, db)
}

func RemoveDayEntry(db *sql.DB, dayEntryId *int) {
	_, err := db.Exec("DELETE FROM dayEntry where id=?", dayEntryId)
	exitOnErr(err, db)
}

func EditDayEntry(db *sql.DB, id int, startTime time.Time, endTime time.Time) {
	_, err := db.Exec("UPDATE dayEntry SET start=?, end=? WHERE id=?", startTime, endTime, id)
	exitOnErr(err, db)
}

func exitOnErr(err error, db *sql.DB) {
	if err != nil {
		log.Fatal(err)
		db.Close()
		os.Exit(1)
	}
}
