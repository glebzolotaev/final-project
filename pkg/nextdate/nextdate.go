package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

const maxSearchDays = 36600

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("repetition rule is not specified")
	}

	date, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("incorrect date %q: %w", dstart, err)
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	parts := strings.Fields(repeat)
	rule := parts[0]
	args := parts[1:]

	switch rule {
	case "d":
		return nextByDays(date, now, args)
	case "y":
		return nextByYears(date, now, args)
	case "w":
		return nextByWeekdays(date, now, args)
	case "m":
		return nextByMonthDays(date, now, args)
	default:
		return "", fmt.Errorf("unsupported format: %q", rule)
	}
}

func nextByDays(date, now time.Time, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("rule d: interval in days is not specified")
	}
	days, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("rule d: incorrect interval %q", args[0])
	}
	if days < 1 {
		return "", errors.New("rule d:  interval must be a positive number")
	}
	if days > 400 {
		return "", errors.New("rule d: interval cannot exceed 400 days")
	}
	for {
		date = date.AddDate(0, 0, days)
		if date.After(now) {
			return date.Format(DateLayout), nil
		}
	}
}

func nextByYears(date, now time.Time, args []string) (string, error) {
	if len(args) > 0 {
		return "", errors.New("rule does not accept parameters")
	}
	for {
		date = date.AddDate(1, 0, 0)
		if date.After(now) {
			return date.Format(DateLayout), nil
		}
	}
}

func nextByWeekdays(date, now time.Time, args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", errors.New("the days of the week are not specified")
	}
	if len(args) > 1 {
		return "", errors.New("too many parameters")
	}

	allowed := [8]bool{}
	for _, s := range strings.Split(args[0], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 7 {
			return "", fmt.Errorf("invalid value %q", s)
		}
		allowed[n] = true
	}

	return search(date, now, func(d time.Time) bool {
		return allowed[weekdayNumber(d)]
	})
}

func weekdayNumber(t time.Time) int {
	return (int(t.Weekday())+6)%7 + 1
}

func nextByMonthDays(date, now time.Time, args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", errors.New("the days of the month are not specified")
	}
	if len(args) > 2 {
		return "", errors.New("too many parameters")
	}

	daySet := map[int]bool{}
	for _, s := range strings.Split(args[0], ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return "", fmt.Errorf("unacceptable day %q", s)
		}
		if n == -1 || n == -2 {
			daySet[n] = true
			continue
		}
		if n < 1 || n > 31 {
			return "", fmt.Errorf("invalid day of the month %d", n)
		}
		daySet[n] = true
	}

	monthSet := map[int]bool{}
	if len(args) == 2 {
		for _, s := range strings.Split(args[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 12 {
				return "", fmt.Errorf("unacceptable month %q", s)
			}
			monthSet[n] = true
		}
	}

	return search(date, now, func(d time.Time) bool {
		if len(monthSet) > 0 && !monthSet[int(d.Month())] {
			return false
		}
		day := d.Day()
		if daySet[day] {
			return true
		}

		last := lastDayOfMonth(d)
		if daySet[-1] && day == last {
			return true
		}
		return daySet[-2] && day == last-1
	})
}

func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

func search(date, now time.Time, matches func(time.Time) bool) (string, error) {
	for i := 0; i < maxSearchDays; i++ {
		date = date.AddDate(0, 0, 1)
		if date.After(now) && matches(date) {
			return date.Format(DateLayout), nil
		}
	}
	return "", errors.New("it was not possible to calculate the next date")
}
