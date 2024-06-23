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

const enterDurationPrompt string = "Please enter a duration in hh:mm or press enter to leave empty\n> "
const editDurationPrompt string = "Please enter a duration in hh:mm or press enter to keep the current value\n> "
const enterTicKeyPrompt string = "Please enter the name of the key for this time entry\n> "
const editTicKeyPrompt string = "Please enter the name of the key for this time entry or press enter to keep the current value\n> "
const enterCommentPrompt string = "Please enter a comment for this time entry\n> "
const editCommentPrompt string = "Please enter a comment for this time entry or press enter to keep the current value\n> "

func AddTimeEntry() (*time.Time, *model.TimeEntry, error) {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Print(enterDatePrompt)
	dateEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, err
	}
	dateResult, err := parseDate(dateEntry)
	if err != nil {
		return nil, nil, err
	}

	buffer.Reset(os.Stdin)
	fmt.Print(enterDurationPrompt)
	duration, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, err
	}
	durationResult, err := parseDuration(duration)
	if err != nil {
		return nil, nil, err
	}

	buffer.Reset(os.Stdin)
	fmt.Print(enterTicKeyPrompt)
	ticKey, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, err
	}
	ticKeyResult := strings.TrimSuffix(string(ticKey), "\n")
	if ticKeyResult == "" {
		return nil, nil, errors.New("key for time entries has to be set")
	}

	buffer.Reset(os.Stdin)
	fmt.Print(enterCommentPrompt)
	comment, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, err
	}
	commentResult := strings.TrimSuffix(string(comment), "\n")
	if commentResult == "" {
		return nil, nil, errors.New("comment for time entries can not be empty")
	}

	resultTimeEntry := model.TimeEntry{
		DayEntryId: 0,
		Duration:   durationResult,
		TicKey:     ticKeyResult,
		Comment:    commentResult,
	}
	return &dateResult, &resultTimeEntry, nil
}

func GetTimeEntry() (*time.Time, string, error) {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Print(enterDatePrompt)
	dateEntry, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, "", err
	}
	dateResult, err := parseDate(dateEntry)
	if err != nil {
		return nil, "", err
	}

	buffer.Reset(os.Stdin)
	fmt.Print(enterTicKeyPrompt)
	ticKey, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, "", err
	}
	ticKeyResult := strings.TrimSuffix(string(ticKey), "\n")
	if ticKeyResult == "" {
		return nil, "", errors.New("key for time entries has to be set")
	}

	return &dateResult, ticKeyResult, nil
}

func EditTimeEntry(timeEntryToEdit *model.TimeEntry) (*float64, *string, *string, error) {
	buffer := bufio.NewReader(os.Stdin)
	fmt.Printf("The duration for the selected time entry is %v\n", timeEntryToEdit.Duration)
	fmt.Print(editDurationPrompt)
	duration, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, nil, err
	}
	var durationResult float64
	if (strings.TrimSuffix(string(duration), "\n")) == "" {
		durationResult = timeEntryToEdit.Duration
	} else {
		durationResult, err = parseDuration(duration)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	buffer.Reset(os.Stdin)
	fmt.Printf("The key for the selected time entry is %v\n", timeEntryToEdit.TicKey)
	fmt.Print(editTicKeyPrompt)
	ticKey, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, nil, err
	}
	ticKeyResult := strings.TrimSuffix(string(ticKey), "\n")
	if ticKeyResult == "" {
		ticKeyResult = timeEntryToEdit.TicKey
	}

	buffer.Reset(os.Stdin)
	fmt.Printf("The comment for the selected time entry is:\n %v\n", timeEntryToEdit.Comment)
	fmt.Print(editCommentPrompt)
	comment, err := buffer.ReadBytes('\n')
	if err != nil {
		return nil, nil, nil, err
	}
	commentResult := strings.TrimSuffix(string(comment), "\n")
	if commentResult == "" {
		commentResult = timeEntryToEdit.Comment
	}
	return &durationResult, &ticKeyResult, &commentResult, nil
}

func parseDuration(duration []byte) (float64, error) {
	var result float64
	var timeRegex string = "^(0?[0-9]|1[0-9]|2[0-3]):([0-5][0-9])"
	matchTimeRegex, err := regexp.Match(timeRegex, duration)
	if err != nil {
		return result, err
	}
	trimmedDuration := strings.TrimSuffix(string(duration), "\n")
	if trimmedDuration == "" {
		return result, nil
	}
	if !matchTimeRegex {
		return result, errors.New("provided did not match the requested format hh:mm")
	}

	regexObj, err := regexp.Compile(":")
	if err != nil {
		return result, err
	}
	timeGroup := regexObj.Split(string(trimmedDuration), 2)
	result, err = constructTimeInHours(timeGroup[0], timeGroup[1])

	return result, err
}

func constructTimeInHours(hour, minute string) (float64, error) {
	var result float64
	hours, err := strconv.Atoi(hour)
	if err != nil {
		return result, err
	}
	minutes, err := strconv.Atoi(minute)
	if err != nil {
		return result, err
	}
	convertedMinutes := float64(minutes) / 60.0
	result = float64(hours) + convertedMinutes
	return result, nil
}
