package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
)

type Config struct {
	Path string `json:"path"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg, err := loadConfig()
	if err != nil {
		fmt.Println("Warning: couldn't load saved config:", err)
		cfg = &Config{}
	}

	var filePath string
	if cfg.Path != "" {
		fmt.Printf("Using saved .ics path: %s\n", cfg.Path)
		fmt.Println("Press Enter to keep it, or type a new path:")
		scanner.Scan()
		if input := scanner.Text(); input != "" {
			filePath = input
		} else {
			filePath = cfg.Path
		}
	} else {
		fmt.Println("Please tell us the path for your .ics file.(If you don't have one create one then tell us the path.)")
		fmt.Println("\n Example: /home/user/Documents/adenda.ics")
		scanner.Scan()
		filePath = scanner.Text()
	}

	if filePath != cfg.Path {
		cfg.Path = filePath
		if err := saveConfig(cfg); err != nil {
			fmt.Println("Warning: couldn't save path for next time:", err)
		}
	}

	fmt.Printf("Your File path is: %v", filePath)
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read ics file %v", err)
	}

	cal, err := ics.ParseCalendar(bytes.NewReader(fileData))
	if err != nil {
		log.Fatalf("Error parsing calendar string: %v\n", err)
	}

	fmt.Println("\nChoose Name of the event")
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
		Description += line
	}
	fmt.Printf("desc: %v", Description)

	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	fmt.Println("\nEnter start time (e.g., 'tomorrow at 8am', 'next monday 15:00', or '26 Jul 8am'):")
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

	event := cal.AddEvent(fmt.Sprintf("%d@adenda.local", time.Now().UnixNano()))
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

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "adenda")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "config.json"), nil
}

func loadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil // no config yet, not an error
		}
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func saveConfig(cfg *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
