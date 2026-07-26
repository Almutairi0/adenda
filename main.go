package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	//"path/filepath"

	//"log"
	"os"
	//"strings"
	ics "github.com/arran4/golang-ical"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"time"
)

func main() {

	filePath := "/home/darling/Documents/adenda.ics"
	fileData, err := os.ReadFile(filePath)

	if err != nil {
		log.Fatalf("Failed to read ics file %v", err)
	}
	cal, err := ics.ParseCalendar(bytes.NewReader(fileData))
	if err != nil {
		log.Fatalf("Error parsing calendar string: %v\n", err)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Choose Name of the event")
	scanner.Scan()
	eventName := scanner.Text()
	fmt.Printf("event name: %v", eventName)
	var Description string
	fmt.Println("\nEnter your Description  (press Enter on a blank line to stop):")
	for scanner.Scan() {
		line := scanner.Text()

		// If the user hits enter without typing anything, break the loop
		if line == "" {
			break
		}
		Description += line + "\n"
	}
	fmt.Printf("desc: %v", Description)
	/*
		timeLayout := "2006-01-02 15:04"

		fmt.Printf("Enter start time (%s): \n", timeLayout)
		scanner.Scan()
		startTimeStr := scanner.Text()
		startTime, err := time.ParseInLocation(timeLayout, startTimeStr, time.Local)
		if err != nil {
			log.Fatalf("Invalid start time format: %v", err)
		}

		fmt.Printf("Enter end time (%s): \n", timeLayout)
		scanner.Scan()
		endTimeStr := scanner.Text()
		endTime, err := time.ParseInLocation(timeLayout, endTimeStr, time.Local)
		if err != nil {
			log.Fatalf("Invalid end time format: %v", err)
		}
	*/

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	fmt.Println("Enter start time (e.g., 'tomorrow at 8am', 'next monday 15:00', or '26 Jul 8am'):")
	scanner.Scan()
	startTimeStr := scanner.Text()

	// Parse relative to right now
	startResult, err := w.Parse(startTimeStr, time.Now())
	if err != nil || startResult == nil {
		log.Fatalf("Could not understand the start time you typed.")
	}
	startTime := startResult.Time

	// 2. Get End Time
	fmt.Println("Enter end time (e.g., 'at 11am', 'in 2 hours'):")
	scanner.Scan()
	endTimeStr := scanner.Text()

	// Parse relative to the START time, not right now!
	endResult, err := w.Parse(endTimeStr, startTime)
	if err != nil || endResult == nil {
		log.Fatalf("Could not understand the end time you typed.")
	}
	endTime := endResult.Time

	event := cal.AddEvent(fmt.Sprintf("%d-%s@adenda.local", time.Now().UnixNano()))
	event.SetSummary(eventName)
	event.SetDescription(Description)
	event.SetStartAt(startTime)
	event.SetEndAt(endTime)
	outputData := cal.Serialize(ics.WithNewLineWindows)

	// Write back to the original file
	err = os.WriteFile(filePath, []byte(outputData), 0644)
	if err != nil {
		log.Fatalf("Failed to save the updated calendar: %v", err)
	}

	fmt.Println("Event successfully added to adenda.ics!")
}
