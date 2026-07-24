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
	"time"

	ics "github.com/arran4/golang-ical"
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
