package date

import (
	"testing"
	"time"
)

func TestYearMonth(t *testing.T) {
	runBindingTest(t, YearMonthOf, YearMonthOnly, false)

	t.Run("Test Construct", func(t *testing.T) {
		now := time.Now()
		local := YearMonthOf(now)
		if !local.Equal(NewYearMonth(now.Year(), now.Month())) {
			t.Fatalf("Expected constructed time give different result")
		}

		t.Run("Checked", func(t *testing.T) {
			if v, err := NewYearMonthChecked(2020, 2); err != nil || !v.Equal(NewYearMonth(2020, 2)) {
				t.Fatalf("NewYearMonthChecked valid: got %v, %v", v, err)
			}
			if _, err := NewYearMonthChecked(2020, 13); err == nil {
				t.Fatalf("NewYearMonthChecked invalid: expected error, got nil")
			}
		})
	})

	t.Run("Test Access YearMonth Properties", func(t *testing.T) {
		now := time.Now()
		local := YearMonthOf(now)
		if local.Month() != now.Month() {
			t.Fatalf("Month: Expected %d, got %d", now.Month(), local.Month())
		}
		if local.Year() != now.Year() {
			t.Fatalf("Year: Expected %d, got %d", now.Year(), local.Year())
		}
		y, m := local.YearMonth()
		y2, m2 := now.Year(), now.Month()
		if y != y2 || m != m2 {
			t.Fatalf("YearMonth mismatch")
		}
	})

	t.Run("Test YearMonth Operations", func(t *testing.T) {
		now := NewYearMonth(2023, 10)
		nextMonth := NewYearMonth(2023, 11)
		before := NewYearMonth(2019, time.January)

		if sub := nextMonth.Sub(now); sub != 24*time.Hour*31 {
			t.Fatalf("Sub failed: expected 31 * 24 hours, got %v", nextMonth.Sub(now))
		}
		if addMonths := now.AddMonths(1); !addMonths.Equal(nextMonth) {
			t.Fatalf("AddMonths failed")
		}
		if addYearMonth := now.AddYearMonth(0, 1); !addYearMonth.Equal(nextMonth) {
			t.Fatalf("AddYearMonth failed")
		}
		if monthsSince := now.MonthsSince(before); monthsSince != 57 {
			t.Fatalf("Sub failed: expected 57 months, got %v", now.MonthsSince(before))
		}

		if d := now.AtEndOfMonth(); !d.Equal(NewDate(now.Year(), now.Month(), 31)) {
			t.Fatalf("AtEndOfMonth produce wrong Date")
		}
		if d := now.AtStartOfMonth(); !d.Equal(NewDate(now.Year(), now.Month(), 1)) {
			t.Fatalf("AtStartOfMonth produce wrong Date")
		}
		if d := now.AtDay(10); !d.Equal(NewDate(now.Year(), now.Month(), 10)) {
			t.Fatalf("AtDay produce wrong Date")
		}

		local := NewYearMonth(1998, 02)
		if d := local.AtDay(29); !d.Equal(NewDate(local.Year(), local.Month(), 28)) || local.IsValidDay(29) {
			t.Fatalf("AtDay produce wrong Date %v", d)
		}
	})

	t.Run("Test Clamping", func(t *testing.T) {
		if local := NewYearMonth(2020, 13); !local.Equal(NewYearMonth(2020, 12)) {
			t.Fatalf("Clamp mismatch: expected 2020-12, got %v", local)
		}
	})

	t.Run("Test ParseYearMonth", func(t *testing.T) {
		for _, in := range []string{"2023-10-15", "2023-10"} {
			if _, err := ParseYearMonth(in); err != nil {
				t.Fatalf("ParseYearMonth(%q) error %v", in, err)
			}
		}

		for _, in := range []string{"", "2023", "2023-13-01", "2023-10-15T"} {
			if _, err := ParseYearMonth(in); err == nil {
				t.Fatalf("ParseYearMonth(%q) expected error, got nil", in)
			}
		}
	})
}

func TestMonthDay(t *testing.T) {
	runBindingTest(t, MonthDayOf, MonthDayOnly, false)

	t.Run("Test Construct", func(t *testing.T) {
		now := time.Now()
		local := MonthDayOf(now)
		if !local.Equal(NewMonthDay(now.Month(), now.Day())) {
			t.Fatalf("Expected constructed time give different result")
		}

		t.Run("Checked", func(t *testing.T) {
			if v, err := NewMonthDayChecked(now.Month(), now.Day()); err != nil || !v.Equal(NewMonthDay(now.Month(), now.Day())) {
				t.Fatalf("NewMonthDayChecked valid: got %v, %v", v, err)
			}
			if _, err := NewMonthDayChecked(02, 31); err == nil {
				t.Fatalf("NewMonthDayChecked invalid: expected error, got nil")
			}
		})
	})

	t.Run("Test Access MonthDay Properties", func(t *testing.T) {
		now := time.Now()
		local := MonthDayOf(now)
		if local.Month() != now.Month() {
			t.Fatalf("Month: Expected %d, got %d", now.Month(), local.Month())
		}
		if local.Day() != now.Day() {
			t.Fatalf("Day: Expected %d, got %d", now.Day(), local.Day())
		}
		y, m := local.MonthDay()
		y2, m2 := now.Month(), now.Day()
		if y != y2 || m != m2 {
			t.Fatalf("MonthDay mismatch")
		}
	})

	t.Run("Test MonthDay Operations", func(t *testing.T) {
		local := NewMonthDay(02, 29)
		if v := local.AtYear(2000); !v.Equal(NewDate(2000, 2, 29)) || !local.IsValidYear(2000) {
			t.Fatalf("AtYear failed %v", v)
		}
		if v := local.AtYear(1998); !v.Equal(NewDate(1998, 2, 28)) || local.IsValidYear(1998) {
			t.Fatalf("AtYear failed %v", v)
		}
	})

	t.Run("Test Clamping", func(t *testing.T) {
		if local := NewMonthDay(02, 31); !local.Equal(NewMonthDay(02, 29)) {
			t.Fatalf("Clamp mismatch: expected 02-29, got %v", local)
		}
	})

	t.Run("Test ParseMonthDay", func(t *testing.T) {
		for _, in := range []string{"2023-10-15", "10-15"} {
			if _, err := ParseMonthDay(in); err != nil {
				t.Fatalf("ParseMonthDay(%q) error %v", in, err)
			}
		}

		for _, in := range []string{"", "2023", "2023-13-01", "2023-10-15T", "02-31"} {
			if _, err := ParseMonthDay(in); err == nil {
				t.Fatalf("ParseMonthDay(%q) expected error, got nil", in)
			}
		}
	})
}
