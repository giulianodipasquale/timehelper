# timehelper

A simple command-line tool to convert between various date/time formats and Unix timestamps.

## Features
- Parse **"now"** as input.
- Parse **Unix timestamps** in seconds, milliseconds, microseconds, or nanoseconds (with automatic detection).
- Parse **RFC3339 / RFC3339Nano** datetime strings.
- Output a summary of the parsed time across multiple standard formats.

## Usage

```bash
# Help menu
./timehelper --help

# Basic conversions
./timehelper now
./timehelper 1721845200
./timehelper --unit ms 1721845200000
./timehelper 2026-07-23T18:00:00Z
```

## Installation

To build and use the tool locally:

```bash
go build .
./timehelper <input>
```
