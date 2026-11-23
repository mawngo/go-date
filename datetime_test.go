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
		t.Run("Checked", func(t *testing.T) {
			if v, err := NewDateTimeChecked(2020, 2, 29, 23, 59, 59, 0); err != nil || !v.Equal(NewDateTime(2020, 2, 29, 23, 59, 59, 0)) {
				t.Fatalf("NewDateTimeChecked valid: got %v, %v", v, err)
			}
			for _, in := range []struct {
				y, mo, d, h, mi, s, ns int
			}{
				{2020, 13, 1, 0, 0, 0, 0},
				{2019, 2, 29, 0, 0, 0, 0},
				{2020, 1, 1, 24, 0, 0, 0},
				{2020, 1, 1, 0, 60, 0, 0},
				{2020, 1, 1, 0, 0, 60, 0},
				{2020, 1, 1, 0, 0, 0, 1_000_000_000},
			} {
				if _, err := NewDateTimeChecked(in.y, time.Month(in.mo), in.d, in.h, in.mi, in.s, in.ns); err == nil {
					t.Fatalf("NewDateTimeChecked(%v): expected error, got nil", in)
				}
			}
		})
	})

	t.Run("Test DateTime Operations", func(t *testing.T) {
		now := NewDateTime(2023, 10, 15, 12, 10, 15, 0)
		nextDay := NewDateTime(2023, 10, 16, 14, 15, 15, 0)
		before := NewDateTime(2019, time.January, 1, 0, 0, 0, 0)

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

		if dt := now.AtDate(2023, 1, 1); !dt.Equal(NewDateTime(2023, 1, 1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond())) {
			t.Fatalf("AtDate produce wrong DateTime")
		}
		if dt := now.AtClock(5, 15, 20); !dt.Equal(NewDateTime(now.Year(), now.Month(), now.Day(), 5, 15, 20, 0)) {
			t.Fatalf("AtClock produce wrong DateTime")
		}
		if dt := now.AtClockNano(5, 15, 20, 2); !dt.Equal(NewDateTime(now.Year(), now.Month(), now.Day(), 5, 15, 20, 2)) {
			t.Fatalf("AtClockNano produce wrong DateTime")
		}
		if !now.AtStartOfDay().Equal(NewDateTime(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0)) {
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

		t.Run("AtWeekday", func(t *testing.T) {
			base := NewDateTime(2026, 11, 14, 12, 30, 45, 123_456_789) // Saturday
			tests := []struct {
				name       string
				base       DateTime
				weekday    time.Weekday
				offsetDays int
				want       DateTime
			}{
				{"same day", base, time.Saturday, 0, NewDateTime(2026, 11, 14, 12, 30, 45, 123_456_789)},
				{"next occurrence", base, time.Saturday, 1, NewDateTime(2026, 11, 21, 12, 30, 45, 123_456_789)},
				{"previous occurrence", base, time.Saturday, -1, NewDateTime(2026, 11, 7, 12, 30, 45, 123_456_789)},
				{"one week forward", base, time.Saturday, 7, NewDateTime(2026, 11, 21, 12, 30, 45, 123_456_789)},
				{"one week backward", base, time.Saturday, -7, NewDateTime(2026, 11, 7, 12, 30, 45, 123_456_789)},
				{"offset past occurrence", base, time.Saturday, -8, NewDateTime(2026, 10, 31, 12, 30, 45, 123_456_789)},
				{"coming Monday", base, time.Monday, 0, NewDateTime(2026, 11, 16, 12, 30, 45, 123_456_789)},
				{"next Monday", base, time.Monday, 1, NewDateTime(2026, 11, 16, 12, 30, 45, 123_456_789)},
				{"previous Monday", base, time.Monday, -1, NewDateTime(2026, 11, 9, 12, 30, 45, 123_456_789)},
				{"coming Friday", base, time.Friday, 0, NewDateTime(2026, 11, 20, 12, 30, 45, 123_456_789)},
				{"previous Friday", base, time.Friday, -1, NewDateTime(2026, 11, 13, 12, 30, 45, 123_456_789)},
				{"cross month boundary", NewDateTime(2026, 11, 2, 8, 0, 0, 0), time.Friday, -1, NewDateTime(2026, 10, 30, 8, 0, 0, 0)},
				{"cross year boundary", NewDateTime(2026, 1, 1, 23, 59, 59, 999_999_999), time.Wednesday, -1, NewDateTime(2025, 12, 31, 23, 59, 59, 999_999_999)},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					if got := tt.base.AtWeekday(tt.weekday, tt.offsetDays); !got.Equal(tt.want) {
						t.Fatalf("AtWeekday(%v, %d): expected %v, got %v", tt.weekday, tt.offsetDays, tt.want, got)
					}
				})
			}
		})
	})

	t.Run("Test Clamping", func(t *testing.T) {
		if local := NewDateTime(2025, 12, 31, 25, 70, 70, 0); !local.Equal(NewDateTime(2025, 12, 31, 23, 59, 59, 0)) {
			t.Fatalf("Clamp mismatch: expected 2025-12-31 23:59:59, got %v", local)
		}
		if local := NewDate(2020, 2, 31); !local.Equal(NewDate(2020, 2, 29)) {
			t.Fatalf("Clamp mismatch: expected 2020-02-29, got %v", local)
		}
		if local := NewDate(2020, 13, 40); !local.Equal(NewDate(2020, 12, 31)) {
			t.Fatalf("Clamp mismatch: expected 2020-12-31, got %v", local)
		}
	})

	t.Run("Test 12 months Wrapping", func(t *testing.T) {
		if local := NewDateTime(2025, 12, 31, 23, 0, 0, 0).Add(2 * time.Hour); !local.Equal(NewDateTime(2026, 1, 1, 1, 0, 0, 0)) {
			t.Fatalf("Add wrap mismatch: expected 2026-01-01 01:00:00, got %v", local)
		}
		if local := NewDateTime(2025, 12, 31, 10, 0, 0, 0).AddDate(0, 1, 0); !local.Equal(NewDateTime(2026, 1, 31, 10, 0, 0, 0)) {
			t.Fatalf("AddDate wrap mismatch: expected 2026-01-31 10:00:00, got %v", local)
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
