package sample

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	durationRE = regexp.MustCompile(`([+\-]?)([0-9]+)([a-zA-Z]+)`)
)

func rfc3339(input time.Time) string {
	return input.Format(time.RFC3339)
}

func dateModify(duration string, input time.Time) (time.Time, error) {
	d := input

	var length int
	matches := durationRE.FindAllStringSubmatch(duration, -1)
	for _, groups := range matches {
		length += len(groups[0])

		amount := groups[2]
		if groups[1] == "-" {
			amount = "-" + amount
		}
		unit := groups[3]
		switch unit {
		case "ns", "us", "µs", "ms", "s", "m", "h":
			if dur, err := time.ParseDuration(amount + unit); err != nil {
				return time.Time{}, fmt.Errorf("failed parsing duration '%s': %w", duration, err)
			} else {
				d = d.Add(dur)
			}
		case "d", "day", "days":
			if days, err := strconv.Atoi(amount); err != nil {
				return time.Time{}, fmt.Errorf("failed parsing duration '%s': %w", duration, err)
			} else {
				d = d.AddDate(0, 0, days)
			}
		case "M", "mo", "month", "months":
			if months, err := strconv.Atoi(amount); err != nil {
				return time.Time{}, fmt.Errorf("failed parsing duration '%s': %w", duration, err)
			} else {
				d = d.AddDate(0, months, 0)
			}
		case "y", "yr", "yrs", "year", "years":
			if years, err := strconv.Atoi(amount); err != nil {
				return time.Time{}, fmt.Errorf("failed parsing duration '%s': %w", duration, err)
			} else {
				d = d.AddDate(years, 0, 0)
			}
		}
	}

	if length != len(duration) {
		return time.Time{}, fmt.Errorf("failed parsing duration '%s': invalid format (probably trailing non-tokens)", duration)
	}

	return d, nil
}

func timeFromString(name, s string) (time.Time, error) {
	tmpl, err := newTemplate(name, s)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed parsing '%s' template: %w", name, err)
	}
	var res bytes.Buffer
	if err := tmpl.Execute(&res, nil); err != nil {
		return time.Time{}, fmt.Errorf("failed evaluating '%s' template: %w", name, err)
	} else if d, err := time.Parse(time.RFC3339, res.String()); err != nil {
		return time.Time{}, fmt.Errorf("failed parsing result of '%s' template into RFC3339 date (%s): %w", name, res.String(), err)
	} else {
		return d, nil
	}
}
