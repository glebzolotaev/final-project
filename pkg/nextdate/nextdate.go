package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// DateLayout — формат дат, используемый в проекте.
const DateLayout = "20060102"

// maxSearchDays — защита от «невозможных» правил (например, m 31 2 —
// 31 февраля не существует), чтобы поиск не зациклился навсегда.
// ~100 лет запаса хватает с избытком.
const maxSearchDays = 36600

// NextDate вычисляет следующую дату выполнения задачи.
//
//	now    — время, от которого ищется ближайшая дата;
//	dstart — исходная дата в формате 20060102;
//	repeat — правило повторения: d <дни>, y, w <дни недели>, m <дни> [месяцы].
//
// Возвращает дату в формате 20060102, строго большую now.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if strings.TrimSpace(repeat) == "" {
		return "", errors.New("правило повторения не указано")
	}

	date, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата %q: %w", dstart, err)
	}

	// отбрасываем время, работаем только с датами
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
		return "", fmt.Errorf("неподдерживаемый формат правила: %q", rule)
	}
}

// nextByDays обрабатывает правило "d <дни>".
func nextByDays(date, now time.Time, args []string) (string, error) {
	if len(args) == 0 {
		return "", errors.New("правило d: не указан интервал в днях")
	}
	days, err := strconv.Atoi(args[0])
	if err != nil {
		return "", fmt.Errorf("правило d: некорректный интервал %q", args[0])
	}
	if days < 1 {
		return "", errors.New("правило d: интервал должен быть положительным числом")
	}
	if days > 400 {
		return "", errors.New("правило d: интервал не может превышать 400 дней")
	}
	for {
		date = date.AddDate(0, 0, days)
		if date.After(now) {
			return date.Format(DateLayout), nil
		}
	}
}

// nextByYears обрабатывает правило "y".
func nextByYears(date, now time.Time, args []string) (string, error) {
	if len(args) > 0 {
		return "", errors.New("правило y не принимает параметров")
	}
	for {
		date = date.AddDate(1, 0, 0)
		if date.After(now) {
			return date.Format(DateLayout), nil
		}
	}
}

// nextByWeekdays обрабатывает правило "w <дни недели>".
func nextByWeekdays(date, now time.Time, args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", errors.New("правило w: не указаны дни недели")
	}
	if len(args) > 1 {
		return "", errors.New("правило w: слишком много параметров")
	}

	allowed := [8]bool{}
	for _, s := range strings.Split(args[0], ",") {
		n, err := strconv.Atoi(s)
		if err != nil || n < 1 || n > 7 {
			return "", fmt.Errorf("правило w: недопустимое значение %q", s)
		}
		allowed[n] = true
	}

	return search(date, now, func(d time.Time) bool {
		return allowed[weekdayNumber(d)]
	})
}

// weekdayNumber возвращает номер дня недели: 1 — понедельник, 7 — воскресенье.
func weekdayNumber(t time.Time) int {
	return (int(t.Weekday())+6)%7 + 1
}

// nextByMonthDays обрабатывает правило "m <дни> [месяцы]".
func nextByMonthDays(date, now time.Time, args []string) (string, error) {
	if len(args) == 0 || args[0] == "" {
		return "", errors.New("правило m: не указаны дни месяца")
	}
	if len(args) > 2 {
		return "", errors.New("правило m: слишком много параметров")
	}

	daySet := map[int]bool{}
	for _, s := range strings.Split(args[0], ",") {
		n, err := strconv.Atoi(s)
		if err != nil {
			return "", fmt.Errorf("правило m: некорректный день %q", s)
		}
		if n == -1 || n == -2 {
			daySet[n] = true
			continue
		}
		if n < 1 || n > 31 {
			return "", fmt.Errorf("правило m: недопустимый день месяца %d", n)
		}
		daySet[n] = true
	}

	monthSet := map[int]bool{}
	if len(args) == 2 {
		for _, s := range strings.Split(args[1], ",") {
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 12 {
				return "", fmt.Errorf("правило m: недопустимый месяц %q", s)
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
		// обработка -1 и -2 (последний и предпоследний день месяца)
		last := lastDayOfMonth(d)
		if daySet[-1] && day == last {
			return true
		}
		return daySet[-2] && day == last-1
	})
}

// lastDayOfMonth возвращает номер последнего дня месяца (28–31).
func lastDayOfMonth(t time.Time) int {
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, t.Location()).Day()
}

// search перебирает дни, начиная со следующего за dstart, пока не найдёт
// дату, которая одновременно больше now и удовлетворяет правилу.
func search(date, now time.Time, matches func(time.Time) bool) (string, error) {
	for i := 0; i < maxSearchDays; i++ {
		date = date.AddDate(0, 0, 1)
		if date.After(now) && matches(date) {
			return date.Format(DateLayout), nil
		}
	}
	return "", errors.New("не удалось вычислить следующую дату по правилу")
}
