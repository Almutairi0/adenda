# adenda

A small command-line tool for adding events to an existing `.ics` calendar file using natural-language date and time input.

## What it does

`adenda` opens an `.ics` calendar file, walks you through a few prompts in the terminal, and appends a new event to it:

1. Reads and parses your existing `.ics` file
2. Prompts for an event name
3. Prompts for a description (type your description, press Enter on a blank line to finish)
4. Prompts for a start time, understood as natural language (e.g. `tomorrow at 8am`, `next monday 15:00`, `26 Jul 8am`)
5. Prompts for an end time, parsed relative to the start time (e.g. `at 11am`, `in 2 hours`)
6. Adds the event to the calendar and writes the file back to disk

Natural-language parsing is handled by [olebedev/when](https://github.com/olebedev/when), and calendar reading/writing is handled by [arran4/golang-ical](https://github.com/arran4/golang-ical).

## Requirements

- Go 1.26.5 or later

## Installation

```bash
git clone https://github.com/Almutairi0/adenda.git
cd adenda
go build -o adenda
```

## Usage

```bash
./adenda
```

Then follow the interactive prompts to add an event.

> **Note:** the path to the `.ics` file is currently hardcoded in `main.go` (`filePath := "/home/darling/Documents/adenda.ics"`). Update this path to point at your own calendar file before building, or adapt the code to accept it as a command-line argument or environment variable.

## Dependencies

- [github.com/arran4/golang-ical](https://github.com/arran4/golang-ical) — parsing and writing `.ics` calendar files
- [github.com/olebedev/when](https://github.com/olebedev/when) — natural-language date/time parsing
