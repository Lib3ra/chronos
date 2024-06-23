package main

import (
	"bufio"
	"chronos/actions"
	repository "chronos/db"
	"chronos/model"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const dateFormat string = "2006-01-02"

func main() {
	db := repository.InitDB()

	dayEntryFs := flag.NewFlagSet("day", flag.ExitOnError)
	addDayFlag := dayEntryFs.Bool("add", false, "Add a day entry")
	deleteDayFlag := dayEntryFs.Bool("delete", false, "Delete a day entry")
	editDayFlag := dayEntryFs.Bool("edit", false, "Edit the start and end time of a day entry")
	listDaysFlag := dayEntryFs.Bool("list", false, "List all day entries")

	timeEntryFs := flag.NewFlagSet("time", flag.ExitOnError)
	addTimeFlag := timeEntryFs.Bool("add", false, "Add a time entry")
	deleteTimeFlag := timeEntryFs.Bool("delete", false, "Delete a time entry")
	editTimeFlag := timeEntryFs.Bool("edit", false, "Edit a time entry")
	listTimeFlag := timeEntryFs.Bool("list", false, "List all time entries for a given day")

	if len(os.Args) < 2 {
		fmt.Printf("Add a command: \n %v\n %v\n", dayEntryFs.Name(), timeEntryFs.Name())
		os.Exit(0)
	}

	switch os.Args[1] {
	case "day":
		if err := dayEntryFs.Parse(os.Args[2:]); err != nil {
			dayEntryFs.ErrorHandling()
			os.Exit(1)
		}
		if nargs := dayEntryFs.NArg(); nargs > 0 || len(os.Args) < 3 {
			dayEntryFs.PrintDefaults()
			os.Exit(1)
		}

		if *addDayFlag {
			dayEntry, err := actions.AddDay()
			date := dayEntry.Date.Format(dateFormat)
			if err != nil {
				log.Fatal(err)
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
				log.Fatal(err)
			}
			idToDelete, err := repository.GetDayEntryIdByDate(db, *dateToDelete)
			if err != nil {
				log.Fatalf("The entry for %v-could not be found.\n", date)
			}
			rows, err := repository.GetAllTimeEntries(db, *idToDelete)
			if err != nil {
				log.Fatal(err)
			}
			if len(*rows) > 0 {
				buffer := bufio.NewReader(os.Stdin)
				fmt.Print("The selected date has associated time entries. Enter y to continue\n")
				answer, err := buffer.ReadString('\n')
				if err != nil {
					log.Fatal(err)
				}
				if strings.TrimSuffix(answer, "\n") == "y" {
					repository.RemoveDayEntry(db, idToDelete)
					fmt.Printf("The entry for %v has been removed.\n", date)
				}
			} else {
				repository.RemoveDayEntry(db, idToDelete)
				fmt.Printf("The entry for %v has been removed.\n", date)
			}
		}
		if *editDayFlag {
			dateToEdit, err := actions.GetDay()
			date := dateToEdit.Format(dateFormat)
			if err != nil {
				log.Fatal(err)
			}
			idToEdit, dayEntryToEdit, err := repository.GetDayEntryByDate(db, *dateToEdit)
			if err != nil {
				log.Fatalf("The entry for %v could not be found.\n", date)
			}
			startTime, endTime, err := actions.EditDay(dayEntryToEdit)
			if err != nil {
				log.Fatal(err)
			}
			repository.EditDayEntry(db, *idToEdit, *startTime, *endTime)
			fmt.Printf("The entry for %v has been updated.\n", date)
		}
		if *listDaysFlag {
			rows, err := repository.GetAllDayEntries(db)
			if err != nil {
				log.Fatal("There was an error listing all day entries\n")
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
	case "time":
		if err := timeEntryFs.Parse(os.Args[2:]); err != nil {
			timeEntryFs.ErrorHandling()
			os.Exit(1)
		}
		if nargs := timeEntryFs.NArg(); nargs > 0 || len(os.Args) < 3 {
			timeEntryFs.PrintDefaults()
			os.Exit(1)
		}
		if *addTimeFlag {
			referenceDate, timeEntry, err := actions.AddTimeEntry()
			if err != nil {
				log.Fatal(err)
			}
			dayEntryId, err := repository.GetDayEntryIdByDate(db, *referenceDate)
			if err != nil {
				fmt.Printf("No entry for %v exists.\n", referenceDate)
			}
			if dayEntryId != nil {
				timeEntryId, err := repository.GetTimeEntryIdByKey(db, *dayEntryId, timeEntry.TicKey)
				if timeEntryId != nil {
					log.Fatalf("Time entry for Key: %v already exists.\n", timeEntry.TicKey)
				}
				if err != nil {
					timeEntry = &model.TimeEntry{
						DayEntryId: *dayEntryId,
						Duration:   timeEntry.Duration,
						TicKey:     timeEntry.TicKey,
						Comment:    timeEntry.Comment,
					}
					repository.InsertTimeEntry(db, timeEntry)
					fmt.Printf("A time entry for %v has been created.\n", referenceDate)
				}
			}
		}
		if *deleteTimeFlag {
			referenceDate, ticKey, err := actions.GetTimeEntry()
			if err != nil {
				log.Fatal(err)
			}
			referenceDateId, err := repository.GetDayEntryIdByDate(db, *referenceDate)
			if err != nil {
				log.Fatal(err)
			}
			timeEntryId, err := repository.GetTimeEntryIdByKey(db, *referenceDateId, ticKey)
			if err != nil {
				log.Fatalf("No key with value %v has been found for %v\n", ticKey,
					referenceDate.Format(dateFormat))
			}
			repository.RemoveTimeEntry(db, *timeEntryId)
			fmt.Printf("The entry for %v has been removed.\n", ticKey)
		}
		if *editTimeFlag {
			referenceDate, ticKey, err := actions.GetTimeEntry()
			if err != nil {
				log.Fatal(err)
			}
			referenceDateId, err := repository.GetDayEntryIdByDate(db, *referenceDate)
			if err != nil {
				log.Fatal(err)
			}
			timeEntryId, timeEntry, err := repository.GetTimeEntryByKey(db, *referenceDateId, ticKey)
			if err != nil {
				fmt.Println(err)
				log.Fatalf("No key with value %v has been found for %v\n", ticKey,
					referenceDate.Format(dateFormat))
			}
			newDuration, newTicKey, newComment, err := actions.EditTimeEntry(timeEntry)
			if err != nil {
				log.Fatal(err)
			}
			repository.EditTimeEntry(db, *timeEntryId, *newDuration, *newTicKey, *newComment)
			fmt.Printf("The entry for %v has been updated.\n", ticKey)
		}
		if *listTimeFlag {
			date, err := actions.GetDay()
			if err != nil {
				log.Fatal(err)
			}
			dateEntryId, err := repository.GetDayEntryIdByDate(db, *date)
			if err != nil {
				log.Fatal(err)
			}
			rows, err := repository.GetAllTimeEntries(db, *dateEntryId)
			if err != nil {
				log.Fatal("There was an error listing all time entries\n")
			}
			fmt.Printf("Date       | Duration | TicKey  | Comment\n")
			for _, v := range *rows {
				fmt.Printf("%v | %v  | %v | %v |\n",
					date.Format(dateFormat), v.Duration, v.TicKey, v.Comment)
			}
		}
	default:
		fmt.Printf("Add a command: \n %v\n %v\n", dayEntryFs.Name(), timeEntryFs.Name())
		os.Exit(0)
	}
	db.Close()
}
