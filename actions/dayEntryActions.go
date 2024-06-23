package actions

import (
	"bufio"
	"chronos/model"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const enterDatePrompt string = "Please enter a date in yyyy-mm-dd format or press enter for current date \n> "
const enterStartTimePrompt string = "Please enter:\n - a start time in hh:mm format or press\n - none\n - press enter to use current time\n> "
const enterEndTimePrompt string = "Please enter:\n - an end time in hh:mm format or press\n - none\n - press enter to use current time\n> "

func AddDay() (*model.DayEntry, error) {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Print(enterDatePrompt)
	dateEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	dateResult, err := parseDate(dateEntry)
	if err != nil {
		return nil, err
	}
	// Read second argument
	buffer.Reset(os.Stdin)
	fmt.Print(enterStartTimePrompt)
	startTimeEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	startTimeResult, err := parseTime(startTimeEntry, dateResult)
	if err != nil {
		return nil, err
	}
	// Read third argument
	buffer.Reset(os.Stdin)
	fmt.Print(enterEndTimePrompt)
	endTimeEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	endTimeResult, err := parseTime(endTimeEntry, dateResult)
	if err != nil {
		return nil, err
	}
	resultEntry := model.DayEntry{
		Date:  dateResult,
		Start: startTimeResult,
		End:   endTimeResult,
	}
	return &resultEntry, nil
}

func DeleteDay() (*time.Time, error) {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Print(enterDatePrompt)
	dateEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	dateResult, err := parseDate(dateEntry)
	if err != nil {
		return nil, err
	}
	return &dateResult, nil
}

func GetDay() (*time.Time, error) {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Print(enterDatePrompt)
	dateEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	dateResult, err := parseDate(dateEntry)
	if err != nil {
		return nil, err
	}
	return &dateResult, err
}
func EditDay(dayEntryToEdit *model.DayEntry) (*time.Time, *time.Time, error) {
	buffer := bufio.NewReader(os.Stdin)
	startHour, startMinutes, _ := dayEntryToEdit.Start.Clock()
	fmt.Printf("The start time for the selected day is %v:%v\n", startHour, startMinutes)
	fmt.Print(enterStartTimePrompt)
	startTimeEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, err
	}
	startTimeResult, err := parseTime(startTimeEntry, dayEntryToEdit.Date)
	if err != nil {
		return nil, nil, err
	}
	if !isValidUpdateScenario(&dayEntryToEdit.Start, &startTimeResult) {
		startTimeResult = dayEntryToEdit.Start
	}

	buffer.Reset(os.Stdin)
	endHour, endMinutes, _ := dayEntryToEdit.End.Clock()
	fmt.Printf("The end time for the selected day is %v:%v\n", endHour, endMinutes)
	fmt.Print(enterEndTimePrompt)
	endTimeEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, err
	}
	endTimeResult, err := parseTime(endTimeEntry, dayEntryToEdit.Date)
	if err != nil {
		return nil, nil, err
	}
	if !isValidUpdateScenario(&dayEntryToEdit.End, &endTimeResult) {
		endTimeResult = dayEntryToEdit.End
	}

	return &startTimeResult, &endTimeResult, err
}

func parseDate(date []byte) (time.Time, error) {
	var result time.Time
	var dateRegex string = "^\\d{4}-(0?[1-9]|1[012])-(0?[1-9]|[12][0-9]|3[01])"
	matchDateFormat, err := regexp.Match(dateRegex, date)
	if err != nil {
		return result, err
	}

	trimmedDate := strings.TrimSuffix(string(date), "\n")
	if trimmedDate == "" {
		return time.Date(
			time.Now().Local().Year(),
			time.Now().Local().Month(),
			time.Now().Local().Day(),
			0, 0, 0, 0, time.Local), nil
	}

	if !matchDateFormat {
		return result, errors.New("provided date did not match the requested format yyyy-mm-dd")
	}

	regexObj, err := regexp.Compile("-")
	if err != nil {
		return result, err
	}
	datesGroup := regexObj.Split(trimmedDate, 3)
	result, err = constructDate(datesGroup[0], datesGroup[1], datesGroup[2])

	return result, err

}

func constructDate(year, month, day string) (time.Time, error) {
	var result time.Time

	years, err := strconv.Atoi(year)
	if err != nil {
		return result, err
	}
	months, err := strconv.Atoi(month)
	if err != nil {
		return result, err
	}
	days, err := strconv.Atoi(day)
	if err != nil {
		return result, err
	}
	result = time.Date(years, time.Month(months), days, 0, 0, 0, 0, time.Local)
	return result, nil
}

func parseTime(timeString []byte, referenceDate time.Time) (time.Time, error) {
	var result time.Time
	var timeRegex string = "^(0?[0-9]|1[0-9]|2[0-3]):([0-5][0-9])"
	matchTimeRegex, err := regexp.Match(timeRegex, timeString)
	if err != nil {
		return result, err
	}
	trimmedTimeString := strings.TrimSuffix(string(timeString), "\n")
	if trimmedTimeString == "none" {
		return result, nil
	}
	if trimmedTimeString == "" {
		result = time.Now().Local()
		return result, nil
	}
	if !matchTimeRegex {
		return result, errors.New("provided did not match the requested format hh:mm")
	}

	regexObj, err := regexp.Compile(":")
	if err != nil {
		return result, err
	}
	timeGroup := regexObj.Split(string(trimmedTimeString), 2)
	result, err = constructTime(timeGroup[0], timeGroup[1], referenceDate)

	return result, err
}

func constructTime(hour, minute string, referenceDate time.Time) (time.Time, error) {
	var result time.Time

	hours, err := strconv.Atoi(hour)
	if err != nil {
		return result, err
	}
	minutes, err := strconv.Atoi(minute)
	if err != nil {
		return result, err
	}
	result = referenceDate.Local().Add(time.Hour*time.Duration(hours) +
		time.Minute*time.Duration(minutes))

	return result, nil
}

func isValidUpdateScenario(before *time.Time, after *time.Time) bool {
	if !before.IsZero() && after.IsZero() {
		return false
	}
	return true
}
