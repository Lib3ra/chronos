package repository

import (
	"chronos/model"
	"database/sql"
)

func InsertTimeEntry(db *sql.DB, timeEntry *model.TimeEntry) {
	_, err := db.Exec("INSERT INTO timeEntry VALUES(null,?,?,?,?)",
		timeEntry.DayEntryId,
		timeEntry.Duration,
		timeEntry.TicKey,
		timeEntry.Comment)
	exitOnErr(err, db)
}

func RemoveTimeEntry(db *sql.DB, timeEntryId int) {
	_, err := db.Exec("DELETE FROM timeEntry where id=?", timeEntryId)
	exitOnErr(err, db)
}

func GetTimeEntryIdByKey(db *sql.DB, referenceId int, key string) (*int, error) {
	var resultId int
	row := db.QueryRow("SELECT id from timeEntry where dayEntryId=? and ticKey=?", referenceId, key)
	err := row.Scan(&resultId)
	if err != nil {
		return nil, err
	}
	return &resultId, nil
}

func GetTimeEntryByKey(db *sql.DB, referenceId int, key string) (*int, *model.TimeEntry, error) {
	var resultId int
	var result model.TimeEntry
	row := db.QueryRow("SELECT id, duration, ticKey, comment from timeEntry where dayEntryId=? and ticKey=?", referenceId, key)
	err := row.Scan(&resultId, &result.Duration, &result.TicKey, &result.Comment)
	if err != nil {
		return nil, nil, err
	}
	return &resultId, &result, nil
}

func EditTimeEntry(db *sql.DB, id int, duration float64, ticKey, comment string) {
	_, err := db.Exec("UPDATE timeEntry SET duration=?, ticKey=?, comment=? WHERE id=?",
		duration, ticKey, comment, id)
	exitOnErr(err, db)
}
func GetAllTimeEntries(db *sql.DB, dayEntryId int) (*[]model.TimeEntry, error) {
	rows, err := db.Query("SELECT duration, ticKey, comment from timeEntry WHERE dayEntryId=? ORDER BY id DESC LIMIT 100", dayEntryId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []model.TimeEntry{}
	for rows.Next() {
		i := model.TimeEntry{}
		err := rows.Scan(&i.Duration, &i.TicKey, &i.Comment)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}
	return &result, nil
}
