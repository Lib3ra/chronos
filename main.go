package main

import (
	"chronos/actions"
	repository "chronos/db"
	"flag"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

/*
Done
create sqlite on startup if not already existing
check if db already created
Move everything db related to different file
abstract away db operations
create datamodel for dayentry
  - date
  - start time
  - end time

add day entry
delete day entry
edit day entry
list day entries

Todo
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

design better control flow
  - dont make user restart from scratch
  - handle errors without necessarily crashing

introduce a logging framework

make -listDays not scuffed
visual presentation when editing
  - navigate with arrowkeys/hjkl
  - try something like the add dialog
*/
const dateFormat string = "2006-01-02"

func main() {
	// DB Conection code
	db := repository.InitDB()
	// FLAG CODE
	addDayFlag := flag.Bool("addDay", false, "Add a day entry")
	deleteDayFlag := flag.Bool("deleteDay", false, "Delete a day entry")
	editDayFlag := flag.Bool("editDay", false, "Edit the start and end time of a day entry")
	listDaysFlag := flag.Bool("listDays", false, "List all day entries")
	flag.Parse()
	if *addDayFlag {
		dayEntry, err := actions.AddDay()
		date := dayEntry.Date.Format(dateFormat)
		if err != nil {
			fmt.Println(err)
		}
		id, err := repository.GetDayEntryIdByDate(db, dayEntry.Date)
		if id != nil {
			fmt.Printf("An entry for %v already exists.\n", date)
		}
		if err != nil {
			repository.InsertDayEntry(db, dayEntry)
			fmt.Printf("An entry for %v has been created.\n", date)
		}
	}
	if *deleteDayFlag {
		dateToDelete, err := actions.DeleteDay()
		date := dateToDelete.Format(dateFormat)
		if err != nil {
			fmt.Println(err)
		}
		idToDelete, err := repository.GetDayEntryIdByDate(db, *dateToDelete)
		if err != nil {
			fmt.Printf("The entry for %v-could not be found.\n", date)
		}
		repository.RemoveDayEntry(db, idToDelete)
		fmt.Printf("The entry for %v- has been removed.\n", date)
	}
	if *editDayFlag {
		dateToEdit, err := actions.GetDay()
		date := dateToEdit.Format(dateFormat)
		if err != nil {
			fmt.Println(err)
		}
		idToEdit, dayEntryToEdit, err := repository.GetDayEntryByDate(db, *dateToEdit)
		if err != nil {
			fmt.Printf("The entry for %v could not be found.\n", date)
		}
		startTime, endTime, err := actions.EditDay(dayEntryToEdit)
		if err != nil {
			fmt.Println(err)
		}
		repository.EditDayEntry(db, *idToEdit, *startTime, *endTime)
		fmt.Printf("The entry for %v has been updated.\n", date)
	}
	if *listDaysFlag {
		rows, err := repository.GetAllDayEntries(db)
		if err != nil {
			fmt.Printf("There was an error listing all day entries\n")
		}
		fmt.Printf("Date       | Start | End  |\n")
		for _, v := range *rows {
			date := v.Date.Format(dateFormat)
			startHours, startMinutes, _ := v.Start.Clock()
			endHours, endMinutes, _ := v.Start.Clock()
			fmt.Printf("%v | %v:%v  | %v:%v |\n",
				date, startHours, startMinutes, endHours, endMinutes)
		}
	}

	db.Close()
}
