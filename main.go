package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/tj/go-naturaldate"
	"github.com/urfave/cli/v2"
)

const hpfDateTimeUTC = "20060102150405.000000"
const description = `timehelper: Convert between date/time formats and Unix timestamps

Usage:
  timehelper [options] <input>
  timehelper [options] <unix-timestamp>
  timehelper [options] <datetime-string>


Features:
  - Parse natural language dates like "now", "tomorrow", "next week", "last week", "a year ago"
  - Parse Unix timestamps in seconds, milliseconds, microseconds, or nanoseconds (with automatic detection)
  - Parse RFC3339 / RFC3339Nano datetime strings
  - Output a summary of the parsed time across multiple standard formats

Examples:
  timehelper now
  timehelper tomorrow
  timehelper "next week"
  timehelper "a year ago"
  timehelper 1721845200
  timehelper --unit ms 1721845200000
  timehelper 2026-07-23T18:00:00Z`

var version = "dev"

func main() {
	app := newApp()

	if err := app.Run(os.Args); err != nil {
		code := 1
		if ec, ok := err.(cli.ExitCoder); ok {
			code = ec.ExitCode()
		}

		fmt.Fprintln(os.Stderr, err)
		os.Exit(code)
	}
}

func newApp() *cli.App {
	return &cli.App{
		Name:                 "timehelper",
		Usage:                "Convert between date/time formats and Unix timestamps",
		UsageText:            "timehelper [options] <input>",
		Description:          description,
		Version:              version,
		EnableBashCompletion: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "unit",
				Aliases: []string{"u"},
				Value:   "auto",
				Usage:   "Unix timestamp unit: auto|s|ms|us|ns",
			},
		},
		Action: run,
	}
}

func run(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return cli.Exit(
			`expected exactly 1 argument; run with --help for usage examples`,
			2,
		)
	}

	input := ctx.Args().First()
	unit := ctx.String("unit")

	t, err := parseInput(input, unit)
	if err != nil {
		return cli.Exit(err.Error(), 2)
	}

	printFormats(os.Stdout, t)
	return nil
}

func parseInput(input string, unit string) (time.Time, error) {
	// Handle Unix timestamps (in seconds, milliseconds, microseconds, or nanoseconds)
	if n, err := strconv.ParseInt(input, 10, 64); err == nil {
		return parseUnix(n, input, unit)
	}

	// Try parsing with naturaldate first (for phrases like "tomorrow", "next week", etc.)
	if t, err := naturaldate.Parse(input, time.Now()); err == nil {
		return t, nil
	}

	// Try parsing as a datetime string (RFC3339/RFC3339Nano)
	if t, err := parseDateTime(input); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf(
		"unsupported input %q; expected natural language date (\"now\", \"tomorrow\", \"next week\", etc...), a Unix timestamp, "+
			"or a RFC3339/RFC3339Nano datetime string",
		input,
	)
}

func parseUnix(value int64, raw string, unit string) (time.Time, error) {
	switch strings.ToLower(unit) {
	case "s":
		return time.Unix(value, 0), nil
	case "ms":
		return time.UnixMilli(value), nil
	case "us":
		return time.UnixMicro(value), nil
	case "ns":
		return time.Unix(0, value), nil
	case "auto":
		digits := len(strings.TrimLeft(raw, "+-"))

		switch digits {
		case 10:
			return time.Unix(value, 0), nil
		case 13:
			return time.UnixMilli(value), nil
		case 16:
			return time.UnixMicro(value), nil
		case 19:
			return time.Unix(0, value), nil
		default:
			return time.Time{}, fmt.Errorf(
				"ambiguous Unix timestamp %q; use --unit s|ms|us|ns",
				raw,
			)
		}
	default:
		return time.Time{}, fmt.Errorf(
			"invalid --unit %q; expected auto|s|ms|us|ns",
			unit,
		)
	}
}

func parseDateTime(input string) (time.Time, error) {
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, input); err == nil {
			return t, nil
		}
	}

	if t, err := time.ParseInLocation(hpfDateTimeUTC, input, time.UTC); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("unsupported datetime format")
}

func printFormats(w io.Writer, t time.Time) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)

	fmt.Fprintf(tw, "Unix:\t%d\n", t.Unix())
	fmt.Fprintf(tw, "Unix Milli:\t%d\n", t.UnixMilli())
	fmt.Fprintf(tw, "Unix Micro:\t%d\n", t.UnixMicro())
	fmt.Fprintf(tw, "Unix Nano:\t%d\n", t.UnixNano())
	fmt.Fprintf(tw, "RFC3339Nano:\t%s\n", t.Format(time.RFC3339Nano))
	fmt.Fprintf(
		tw,
		"RFC3339Nano UTC:\t%s\n",
		t.UTC().Format(time.RFC3339Nano),
	)
	fmt.Fprintf(
		tw,
		"HPFDateTime UTC:\t%s\n",
		t.UTC().Format(hpfDateTimeUTC),
	)

	_ = tw.Flush()
}
