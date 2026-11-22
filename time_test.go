package date

import (
	"strings"
	"testing"
	"time"
)

type timePtr[T any] interface {
	localTime[T]
	*T
}

func runTimeBasicTest[T any, P timePtr[T]](t *testing.T, init func(t time.Time) T) {
	now := time.Now()
	local := init(now)
	t.Run("Test Access Time Properties", func(t *testing.T) {
		local := P(&local)
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
		init := func(t time.Time) P {
			v := init(t)
			return P(&v)
		}
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
		if !local.Equal(NewTimeNano(now.Hour(), now.Minute(), now.Second(), now.Nanosecond())) {
			t.Fatalf("Expected constructed time give different result")
		}

		if NewTime(0, 0, 0).IsZero() {
			t.Fatalf("Expected midnight not to be zero")
		}
		if NewTimeNano(0, 0, 0, 0).IsZero() {
			t.Fatalf("Expected midnight not to be zero")
		}
		if !NewTimeNano(12, 15, 11, 0).Equal(NewTime(12, 15, 11)) {
			t.Fatalf("Expected constructed time give same result")
		}

		t.Run("Checked", func(t *testing.T) {
			if v, err := NewTimeNanoChecked(23, 59, 59, 999_999_999); err != nil || !v.Equal(NewTimeNano(23, 59, 59, 999_999_999)) {
				t.Fatalf("NewTimeNanoChecked valid: got %v, %v", v, err)
			}
			if v, err := NewTimeNanoChecked(12, 30, 45, 0); err != nil || !v.Equal(NewTime(12, 30, 45)) {
				t.Fatalf("NewTimeNanoChecked valid: got %v, %v", v, err)
			}
			for _, in := range [][4]int{
				{24, 0, 0, 0},
				{0, 60, 0, 0},
				{0, 0, 60, 0},
				{0, 0, 0, 1_000_000_000},
			} {
				if _, err := NewTimeNanoChecked(in[0], in[1], in[2], in[3]); err == nil {
					t.Fatalf("NewTimeNanoChecked(%v): expected error, got nil", in)
				}
			}
		})

		func() {
			b, err := NewTime(0, 0, 0).MarshalJSON()
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
		now := NewTime(23, 10, 15)

		if sub := now.Sub(NewTimeNano(19, 0, 15, 0)); sub != (4*time.Hour + 10*time.Minute) {
			t.Fatalf("Sub failed: expected 4 hours, got %v", now.Sub(NewTimeNano(4, 0, 15, 0)))
		}
		if add := now.Add(10 * time.Minute); !add.Equal(NewTimeNano(23, 20, 15, 0)) {
			t.Fatalf("Add failed")
		}

		func() {
			now := time.Now()
			local := TimeOf(now)
			h, m, s, nano := local.ClockNano()
			h2, m2, s2 := now.Clock()
			nano2 := now.Nanosecond()
			if h != h2 || m != m2 || s != s2 || nano != nano2 {
				t.Fatalf("Time mismatch")
			}
		}()

		if dt := now.AtDate(2023, 1, 1); !dt.Equal(NewDateTime(2023, 1, 1, now.Hour(), now.Minute(), now.Second(), now.Nanosecond())) {
			t.Fatalf("AtDate produce wrong DateTime")
		}
	})

	t.Run("Test Clamping", func(t *testing.T) {
		if local := NewTime(25, 70, 70); !local.Equal(NewTime(23, 59, 59)) {
			t.Fatalf("Clamp mismatch: expected 23:59:59, got %v", local)
		}
		if local := NewTimeNano(-1, -1, -1, -1); !local.Equal(NewTimeNano(0, 0, 0, 0)) {
			t.Fatalf("Clamp mismatch: expected 00:00:00, got %v", local)
		}
	})

	t.Run("Test ParseTime", func(t *testing.T) {
		tests := []struct {
			name string
			in   string
			want Time
		}{
			{"no fraction", "12:10:15", NewTime(12, 10, 15)},
			{"with fraction", "12:10:15.123456789", NewTimeNano(12, 10, 15, 123_456_789)},
			{"midnight", "00:00:00", NewTime(0, 0, 0)},
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
