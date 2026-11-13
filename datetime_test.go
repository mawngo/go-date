package date

import (
	"testing"
	"time"
)

func TestDateTime(t *testing.T) {
	runBindingTest(t, DateTimeOf, DateTimeNano, true)
	runDateBasicTest(t, DateTimeOf)
	runTimeBasicTest(t, DateTimeOf)

	t.Run("Test Construct", func(t *testing.T) {
		now := time.Now()
		local := DateTimeOf(now)
		if !local.Equal(NewDateTime(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond())) {
			t.Fatalf("Expected constructed time give different result")
		}
	})

	t.Run("Test DateTime Operations", func(t *testing.T) {
		now := NewDateTime(2023, 10, 15, 12, 10, 15, 0)
		nextDay := NewDateTime(2023, 10, 16, 14, 15, 15, 0)
		before := NewDateTime(2019, time.January, 1, 0, 0, 0, 0)
		after := NewDateTime(2023, 12, 15, 14, 15, 10, 0)

		if !now.Equal(NewDateTime(2023, 10, 15, 12, 10, 15, 0)) || now.Equal(nextDay) {
			t.Fatalf("Equally failed")
		}

		if sub := nextDay.Sub(now); sub != 26*time.Hour+5*time.Minute {
			t.Fatalf("Sub failed: expected 26 hours, got %v", nextDay.Sub(now))
		}
		if add := now.Add(26*time.Hour + 5*time.Minute); !add.Equal(nextDay) {
			t.Fatalf("Add failed")
		}

		if addDays := now.AddDays(1); !addDays.Equal(nextDay.AtClock(now.Clock())) {
			t.Fatalf("AddDays failed")
		}
		if addDate := now.AddDate(0, 0, 1); !addDate.Equal(nextDay.AtClock(now.Clock())) {
			t.Fatalf("AddDate failed")
		}
		if days := now.DaysSince(before); days != 1748 {
			t.Fatalf("Sub failed: expected 1748 days, got %v", now.DaysSince(before))
		}
		if !now.After(before) || now.After(after) || now.After(now) {
			t.Fatalf("After failed")
		}
		if !now.Before(after) || now.Before(before) || now.Before(now) {
			t.Fatalf("Before failed")
		}
		//nolint:gocritic
		if now.Compare(now) != 0 || now.Compare(after) != -1 || now.Compare(before) != 1 {
			t.Fatalf("Compare failed")
		}

		if dt := now.AtDate(2023, 1, 1); !dt.Equal(DateTimeFrom(NewDate(2023, 1, 1), NewClock(now.Clock()))) {
			t.Fatalf("AtDate produce wrong DateTime")
		}
		if dt := now.AtClock(5, 15, 20); !dt.Equal(DateTimeFrom(NewDate(now.Date()), NewClock(5, 15, 20))) {
			t.Fatalf("AtClock produce wrong DateTime")
		}
		if dt := now.AtTime(5, 15, 20, 0); !dt.Equal(DateTimeFrom(NewDate(now.Date()), NewTime(5, 15, 20, 0))) {
			t.Fatalf("AtTime produce wrong DateTime")
		}
		if !now.AtStartOfDay().Equal(DateTimeFrom(NewDate(now.Date()), NewClock(0, 0, 0))) {
			t.Fatalf("AtStartOfDay produce wrong DateTime")
		}

		if date := NewDateTime(2025, 1, 1, 12, 1, 12, 12).UTC(); !date.Equal(time.Date(2025, 1, 1, 12, 1, 12, 12, time.UTC)) {
			t.Fatalf("UTC produce wrong Date")
		}
		if date := NewDateTime(2025, 1, 1, 12, 1, 12, 12).Local(); !date.Equal(time.Date(2025, 1, 1, 12, 1, 12, 12, time.Local)) {
			t.Fatalf("Local produce wrong Date")
		}

		tz, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			t.Fatal(err)
		}
		if date := NewDateTime(2025, 1, 1, 12, 1, 12, 12).In(tz); !date.Equal(time.Date(2025, 1, 1, 12, 1, 12, 12, tz)) {
			t.Fatalf("In produce wrong Date")
		}
	})

	t.Run("Test 24h Wrapping", func(t *testing.T) {
		local := NewDateTime(2025, 12, 31, 25, 70, 70, 0)
		if local.Second() != 10 {
			t.Fatalf("Second wrap mismatch: %d", local.Second())
		}
		if local.Minute() != 11 {
			t.Fatalf("Minute wrap mismatch: %d", local.Minute())
		}
		if local.Hour() != 2 {
			t.Fatalf("Hour wrap mismatch: %d", local.Hour())
		}
		if local.Year() != 2026 {
			t.Fatalf("Year wrap mismatch: %d", local.Year())
		}
		if local.Month() != 1 {
			t.Fatalf("Month wrap mismatch: %d", local.Month())
		}
		if local.Day() != 1 {
			t.Fatalf("Day wrap mismatch: %d", local.Day())
		}
		if !local.Equal(NewDateTime(2026, 1, 1, 2, 11, 10, 0)) {
			t.Fatalf("Equal wrap mismatch - wrapping in constructor is not properly handled")
		}
		if local.Add(23*time.Hour) != NewDateTime(2026, 1, 2, 1, 11, 10, 0) {
			t.Fatalf("Add wrap mismatch")
		}
		if local.Add(-23*time.Hour) != NewDateTime(2025, 12, 31, 3, 11, 10, 0) {
			t.Fatalf("Add wrap mismatch")
		}
	})

	t.Run("Test Truncate/Round", func(t *testing.T) {
		truncateTests := []struct {
			name string
			dt   DateTime
			dur  time.Duration
			want DateTime
		}{
			{"Truncate to second", NewDateTime(2023, 10, 15, 12, 10, 15, 999_999_999), time.Second, NewDateTime(2023, 10, 15, 12, 10, 15, 0)},
			{"Truncate to minute", NewDateTime(2023, 10, 15, 12, 10, 59, 0), time.Minute, NewDateTime(2023, 10, 15, 12, 10, 0, 0)},
			{"Truncate to hour", NewDateTime(2023, 10, 15, 12, 59, 59, 999_999_999), time.Hour, NewDateTime(2023, 10, 15, 12, 0, 0, 0)},
			{"Truncate non-positive unchanged", NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789), 0, NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789)},
		}
		for _, tt := range truncateTests {
			if got := tt.dt.Truncate(tt.dur); !got.Equal(tt.want) {
				t.Fatalf("Truncate failed %s: expected %v, got %v", tt.name, tt.want, got)
			}
		}

		roundTests := []struct {
			name string
			dt   DateTime
			dur  time.Duration
			want DateTime
		}{
			{"Round to second down", NewDateTime(2023, 10, 15, 12, 10, 15, 499_999_999), time.Second, NewDateTime(2023, 10, 15, 12, 10, 15, 0)},
			{"Round to second up", NewDateTime(2023, 10, 15, 12, 10, 15, 500_000_000), time.Second, NewDateTime(2023, 10, 15, 12, 10, 16, 0)},
			{"Round to minute down", NewDateTime(2023, 10, 15, 12, 10, 29, 999_999_999), time.Minute, NewDateTime(2023, 10, 15, 12, 10, 0, 0)},
			{"Round to minute halfway up", NewDateTime(2023, 10, 15, 12, 10, 30, 0), time.Minute, NewDateTime(2023, 10, 15, 12, 11, 0, 0)},
			{"Round to hour halfway up", NewDateTime(2023, 10, 15, 12, 30, 0, 0), time.Hour, NewDateTime(2023, 10, 15, 13, 0, 0, 0)},
			{"Round non-positive unchanged", NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789), 0, NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789)},
		}
		for _, tt := range roundTests {
			if got := tt.dt.Round(tt.dur); !got.Equal(tt.want) {
				t.Fatalf("Round failed %s: expected %v, got %v", tt.name, tt.want, got)
			}
		}
	})

	t.Run("Test ParseDateTime", func(t *testing.T) {
		tests := []struct {
			name string
			in   string
			want DateTime
		}{
			{"space separator", "2023-10-15 12:10:15", NewDateTime(2023, 10, 15, 12, 10, 15, 0)},
			{"space separator with fraction", "2023-10-15 12:10:15.123456789", NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789)},
			{"T separator", "2023-10-15T12:10:15", NewDateTime(2023, 10, 15, 12, 10, 15, 0)},
			{"t separator", "2023-10-15t12:10:15", NewDateTime(2023, 10, 15, 12, 10, 15, 0)},
			{"T separator with fraction", "2023-10-15T12:10:15.123456789", NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789)},
			{"T separator with small fraction", "2023-10-15T12:10:15.12345", NewDateTime(2023, 10, 15, 12, 10, 15, 123_450_000)},
			{"T separator with UTC zone", "2023-10-15T12:10:15Z", NewDateTime(2023, 10, 15, 12, 10, 15, 0)},
			{"t separator with UTC zone", "2023-10-15t12:10:15Z", NewDateTime(2023, 10, 15, 12, 10, 15, 0)},
			{"T separator with offset", "2023-10-15T12:10:15.123456789+07:00", NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789)},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseDateTime(tt.in)
				if err != nil {
					t.Fatalf("ParseDateTime(%q) unexpected error: %v", tt.in, err)
				}
				if !got.Equal(tt.want) {
					t.Fatalf("ParseDateTime(%q): expected %v, got %v", tt.in, tt.want, got)
				}
			})
		}

		for _, in := range []string{"", "2023", "2023-10-15T12:10:15x", "2023-10-15 12:10:15T"} {
			if _, err := ParseDateTime(in); err == nil {
				t.Fatalf("ParseDateTime(%q) expected error, got nil", in)
			}
		}
	})

	t.Run("Test Nanoseconds", func(t *testing.T) {
		if got := NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789).String(); got != "2023-10-15 12:10:15.123456789" {
			t.Fatalf("String mismatch: expected fraction, got %s", got)
		}
		if got := NewDateTime(2023, 10, 15, 12, 10, 15, 120_000_000).String(); got != "2023-10-15 12:10:15.12" {
			t.Fatalf("String mismatch: expected trimmed fraction, got %s", got)
		}
		if got := NewDateTime(2023, 10, 15, 12, 10, 15, 0).String(); got != "2023-10-15 12:10:15" {
			t.Fatalf("String mismatch: expected no fraction, got %s", got)
		}

		dt := NewDateTime(2023, 10, 15, 12, 10, 15, 123_456_789)
		b, err := dt.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		var decoded DateTime
		if err := decoded.UnmarshalJSON(b); err != nil {
			t.Fatal(err)
		}
		if !decoded.Equal(dt) {
			t.Fatalf("JSON round-trip mismatch: expected %v, got %v", dt, decoded)
		}
	})
}
