package date

import (
	"strings"
	"testing"
	"time"
)

func runTimeBasicTest[T localTime[T]](t *testing.T, init func(t time.Time) T) {
	now := time.Now()
	local := init(now)
	t.Run("Test Access Time Properties", func(t *testing.T) {
		if zero := init(time.Time{}); !zero.IsZero() {
			t.Fatalf("Expected zero time, got %v", zero)
		}

		if local.Hour() != now.Hour() {
			t.Fatalf("Hour: Expected %d, got %d", now.Hour(), local.Hour())
		}
		if local.Minute() != now.Minute() {
			t.Fatalf("Minute: Expected %d, got %d", now.Minute(), local.Minute())
		}
		if local.Second() != now.Second() {
			t.Fatalf("Seconds: Expected %d, got %d", now.Second(), local.Second())
		}
		if local.Nanosecond() != now.Nanosecond() {
			t.Fatalf("Nanosecond: Expected %d, got %d", now.Nanosecond(), local.Nanosecond())
		}

		h, m, s := local.Clock()
		h2, m2, s2 := now.Clock()
		if h != h2 || m != m2 || s != s2 {
			t.Fatalf("Clock mismatch")
		}
	})

	t.Run("Test Nanoseconds", func(t *testing.T) {
		now := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 123456789, time.UTC)
		if got := init(now.Truncate(time.Second)).String(); strings.Contains(got, ".") {
			t.Fatalf("String mismatch: expected no fraction, got %s", got)
		}
		if got := init(now.Truncate(time.Microsecond)).String(); strings.Index(got, ".123456") != len(got)-7 {
			t.Fatalf("String mismatch: expected no fraction, got %s", got)
		}
		if got := init(now.Truncate(time.Millisecond)).String(); strings.Index(got, ".123") != len(got)-4 {
			t.Fatalf("String mismatch: expected no fraction, got %s", got)
		}
	})
}

func TestTime(t *testing.T) {
	runBindingTest(t, TimeOf, TimeOnlyNano, false)
	runTimeBasicTest(t, TimeOf)

	t.Run("Test Construct", func(t *testing.T) {
		now := time.Now()
		local := TimeOf(now)
		if !local.Equal(NewTime(now.Hour(), now.Minute(), now.Second(), now.Nanosecond())) {
			t.Fatalf("Expected constructed time give different result")
		}

		if NewClock(0, 0, 0).IsZero() {
			t.Fatalf("Expected midnight not to be zero")
		}
		if NewTime(0, 0, 0, 0).IsZero() {
			t.Fatalf("Expected midnight not to be zero")
		}

		func() {
			b, err := NewClock(0, 0, 0).MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			var d Time
			if err := d.UnmarshalJSON(b); err != nil {
				t.Fatal(err)
			}
			if d.IsZero() || d.String() != "00:00:00" {
				t.Fatalf("Midnight JSON round-trip mismatch: got %v", d)
			}
		}()
	})

	t.Run("Test Time Operations", func(t *testing.T) {
		now := NewClock(23, 10, 15)
		before := NewClock(19, 0, 0)
		after := NewClock(23, 15, 0)

		if !now.Equal(NewTime(23, 10, 15, 0)) || now.Equal(NewTime(23, 10, 15, 1)) {
			t.Fatalf("Equally failed")
		}
		if sub := now.Sub(NewTime(19, 0, 15, 0)); sub != (4*time.Hour + 10*time.Minute) {
			t.Fatalf("Sub failed: expected 4 hours, got %v", now.Sub(NewTime(4, 0, 15, 0)))
		}
		if add := now.Add(10 * time.Minute); !add.Equal(NewTime(23, 20, 15, 0)) {
			t.Fatalf("Add failed")
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

		func() {
			now := time.Now()
			local := TimeOf(now)
			h, m, s, nano := local.Time()
			h2, m2, s2 := now.Clock()
			nano2 := now.Nanosecond()
			if h != h2 || m != m2 || s != s2 || nano != nano2 {
				t.Fatalf("Time mismatch")
			}
		}()

		if dt := now.AtDate(2023, 1, 1); !dt.Equal(DateTimeFrom(NewDate(2023, 1, 1), now)) {
			t.Fatalf("AtDate produce wrong DateTime")
		}
	})

	t.Run("Test 24h Wrapping", func(t *testing.T) {
		local := NewClock(25, 70, 70)
		if local.Second() != 10 {
			t.Fatalf("Second wrap mismatch: %d", local.Second())
		}
		if local.Minute() != 11 {
			t.Fatalf("Minute wrap mismatch: %d", local.Minute())
		}
		if local.Hour() != 2 {
			t.Fatalf("Hour wrap mismatch: %d", local.Hour())
		}
		if !local.Equal(NewClock(2, 11, 10)) {
			t.Fatalf("Equal wrap mismatch - wrapping in constructor is not properly handled")
		}
		if local.Add(23*time.Hour) != NewClock(1, 11, 10) {
			t.Fatalf("Add wrap mismatch")
		}
		if local.Add(-23*time.Hour) != NewClock(3, 11, 10) {
			t.Fatalf("Add wrap mismatch")
		}
	})

	t.Run("Test ParseTime", func(t *testing.T) {
		tests := []struct {
			name string
			in   string
			want Time
		}{
			{"no fraction", "12:10:15", NewClock(12, 10, 15)},
			{"with fraction", "12:10:15.123456789", NewTime(12, 10, 15, 123_456_789)},
			{"midnight", "00:00:00", NewClock(0, 0, 0)},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseTime(tt.in)
				if err != nil {
					t.Fatalf("ParseTime(%q) unexpected error: %v", tt.in, err)
				}
				if !got.Equal(tt.want) {
					t.Fatalf("ParseTime(%q): expected %v, got %v", tt.in, tt.want, got)
				}
			})
		}

		if got, err := ParseTime("00:00:00"); err != nil || got.IsZero() {
			t.Fatalf("ParseTime(00:00:00): expected non-zero midnight, got zero=%v err=%v", got.IsZero(), err)
		}

		for _, in := range []string{"", "12", "12:10", "25:00:00", "12:10:15x"} {
			if _, err := ParseTime(in); err == nil {
				t.Fatalf("ParseTime(%q) expected error, got nil", in)
			}
		}
	})
}
