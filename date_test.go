package date

import (
	"database/sql/driver"
	"fmt"
	"testing"
	"time"
)

type bindingPtr[T any] interface {
	binding
	*T
}

func runBindingTest[T any, P bindingPtr[T]](t *testing.T, init func(time time.Time) T, format string, rfc3339 bool) {
	t.Run("Format/AppendFormat", func(t *testing.T) {
		now := time.Now()
		local := init(now)

		for _, format := range []string{time.RFC3339Nano, time.RFC1123Z, time.RubyDate} {
			expected := P(&local).unwrap().Format(format)
			if P(&local).Format(format) != expected {
				t.Fatalf("Format mismatch: expected %s, got %s", expected, P(&local).Format(format))
			}

			if b := P(&local).AppendFormat(nil, format); expected != string(b) {
				t.Fatalf("AppendFormat mismatch: expected %s, got %s", expected, b)
			}
		}
	})

	t.Run("String/Scan/Value", func(t *testing.T) {
		now := time.Now()
		local := init(now)
		expected := now.Format(format)

		if P(&local).String() != expected {
			t.Fatalf("String mismatch: expected %s, got %s", expected, P(&local).String())
		}
		if v, _ := P(&local).Value(); v != driver.Value(expected) {
			t.Fatalf("Value mismatch: expected %s, got %s", expected, P(&local).String())
		}

		for _, v := range []any{expected, []byte(expected), now} {
			var d T
			if err := P(&d).Scan(v); err != nil {
				t.Fatal(err)
			}
			if P(&d).String() != expected {
				t.Fatalf("Scan mismatch: expected %s, got %s", expected, P(&d).String())
			}
		}

		if rfc3339 {
			t.Run("RFC3339", func(t *testing.T) {
				expected = now.Format(time.RFC3339Nano)
				for _, v := range []any{expected, []byte(expected), now} {
					var d T
					if err := P(&d).Scan(v); err != nil {
						t.Fatal(err)
					}
					if P(&d).String() != now.Format(format) {
						t.Fatalf("Scan mismatch")
					}
				}
			})
		}

		t.Run("Zero", func(t *testing.T) {
			zero := init(time.Time{})
			if P(&zero).String() != (time.Time{}).Format(format) {
				t.Fatalf("Zero String mismatch: expected %s, got %s", time.Time{}.Format(format), P(&zero).String())
			}
			if v, _ := P(&zero).Value(); v != nil {
				t.Fatalf("Zero Value mismatch: expected nil, got %v", v)
			}
			var zeroD T
			if err := P(&zeroD).Scan(nil); err != nil {
				t.Fatal(err)
			}
			if !P(&zeroD).unwrap().IsZero() {
				t.Fatal("Scan nil mismatch: expected zero")
			}
		})
	})

	t.Run("Json Marshal/Unmarshal", func(t *testing.T) {
		now := time.Now()
		expected := fmt.Sprintf("\"%s\"", now.Format(format))

		t.Run("Marshal", func(t *testing.T) {
			local := init(now)
			b, err := P(&local).MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != expected {
				t.Fatalf("MarshalJson mismatch: expected %s, got %s", expected, string(b))
			}
		})

		t.Run("Unmarshal", func(t *testing.T) {
			var d T
			if err := P(&d).UnmarshalJSON([]byte(expected)); err != nil {
				t.Fatal(err)
			}
			if P(&d).String() != now.Format(format) {
				t.Fatal("UnmarshalJson mismatch")
			}
		})

		if rfc3339 {
			t.Run("RFC3339", func(t *testing.T) {
				expected := fmt.Sprintf("\"%s\"", now.Format(time.RFC3339Nano))
				var d T
				if err := P(&d).UnmarshalJSON([]byte(expected)); err != nil {
					t.Fatal(err)
				}
				if P(&d).String() != now.Format(format) {
					t.Fatal("UnmarshalJson mismatch")
				}
			})
		}

		t.Run("Zero", func(t *testing.T) {
			zero := init(time.Time{})
			b, err := P(&zero).MarshalJSON()
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != "null" {
				t.Fatalf("Zero MarshalJson mismatch: expected null, got %s", string(b))
			}
			var zeroD T
			if err := P(&zeroD).UnmarshalJSON([]byte("null")); err != nil {
				t.Fatal(err)
			}
			if !P(&zeroD).unwrap().IsZero() {
				t.Fatal("Zero UnmarshalJson mismatch: expected zero")
			}
		})
	})

	t.Run("Text Marshal/Unmarshal", func(t *testing.T) {
		now := time.Now()
		expected := now.Format(format)

		t.Run("Marshal", func(t *testing.T) {
			local := init(now)
			b, err := P(&local).MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != expected {
				t.Fatalf("MarshalText mismatch: expected %s, got %s", expected, string(b))
			}
		})
		t.Run("Unmarshal", func(t *testing.T) {
			var d T
			if err := P(&d).UnmarshalText([]byte(expected)); err != nil {
				t.Fatal(err)
			}
			if P(&d).String() != expected {
				t.Fatal("UnmarshalText mismatch")
			}
		})

		if rfc3339 {
			t.Run("RFC3339", func(t *testing.T) {
				var d T
				if err := P(&d).UnmarshalText([]byte(now.Format(time.RFC3339Nano))); err != nil {
					t.Fatal(err)
				}
				if P(&d).String() != expected {
					t.Fatal("UnmarshalText mismatch")
				}
			})
		}

		t.Run("Zero", func(t *testing.T) {
			zero := init(time.Time{})
			b, err := P(&zero).MarshalText()
			if err != nil {
				t.Fatal(err)
			}
			if string(b) != "" {
				t.Fatalf("Zero MarshalText mismatch: expected empty, got %s", string(b))
			}
			var zeroD T
			if err := P(&zeroD).UnmarshalText([]byte("")); err != nil {
				t.Fatal(err)
			}
			if !P(&zeroD).unwrap().IsZero() {
				t.Fatal("Zero UnmarshalText mismatch: expected zero")
			}
		})
	})

	t.Run("UnmarshalParam", func(t *testing.T) {
		now := time.Now()
		expected := now.Format(format)
		t.Run("Unmarshal", func(t *testing.T) {
			var d T
			if err := P(&d).UnmarshalParam(expected); err != nil {
				t.Fatal(err)
			}
			if P(&d).String() != expected {
				t.Fatal("UnmarshalParam mismatch")
			}
		})

		if rfc3339 {
			t.Run("RFC3339", func(t *testing.T) {
				var d T
				if err := P(&d).UnmarshalParam(now.Format(time.RFC3339Nano)); err != nil {
					t.Fatal(err)
				}
				if P(&d).String() != expected {
					t.Fatal("UnmarshalParam mismatch")
				}
			})
		}

		t.Run("Zero", func(t *testing.T) {
			var zeroD T
			if err := P(&zeroD).UnmarshalParam(""); err != nil {
				t.Fatal(err)
			}
			if !P(&zeroD).unwrap().IsZero() {
				t.Fatal("Zero UnmarshalParam mismatch: expected zero")
			}
		})
	})

	t.Run("Gob Encode/Decode", func(t *testing.T) {
		now := time.Now()
		local := init(now)

		b, err := P(&local).GobEncode()
		if err != nil {
			t.Fatal(err)
		}

		var d T
		if err := P(&d).GobDecode(b); err != nil {
			t.Fatal(err)
		}
		if P(&d).String() != now.Format(format) {
			t.Fatal("Gob Encode/Decode mismatch")
		}

		t.Run("Zero", func(t *testing.T) {
			zero := init(time.Time{})
			zb, err := P(&zero).GobEncode()
			if err != nil {
				t.Fatal(err)
			}
			var zeroD T
			if err := P(&zeroD).GobDecode(zb); err != nil {
				t.Fatal(err)
			}
			if !P(&zeroD).unwrap().IsZero() {
				t.Fatal("Zero Gob Encode/Decode mismatch: expected zero")
			}
		})
	})

	t.Run("Binary Marshal/Unmarshal", func(t *testing.T) {
		now := time.Now()
		local := init(now)

		b, err := P(&local).MarshalBinary()
		if err != nil {
			t.Fatal(err)
		}

		var d T
		if err := P(&d).UnmarshalBinary(b); err != nil {
			t.Fatal(err)
		}
		if P(&d).String() != now.Format(format) {
			t.Fatal("Binary Marshal/Unmarshal mismatch")
		}

		t.Run("Zero", func(t *testing.T) {
			zero := init(time.Time{})
			zb, err := P(&zero).MarshalBinary()
			if err != nil {
				t.Fatal(err)
			}
			var zeroD T
			if err := P(&zeroD).UnmarshalBinary(zb); err != nil {
				t.Fatal(err)
			}
			if !P(&zeroD).unwrap().IsZero() {
				t.Fatal("Zero Binary Marshal/Unmarshal mismatch: expected zero")
			}
		})
	})
}

func runDateBasicTest[T localDate[T]](t *testing.T, init func(time.Time) T) {
	now := time.Now()
	local := init(now)
	t.Run("Test Access Date Properties", func(t *testing.T) {
		if zero := init(time.Time{}); !zero.IsZero() {
			t.Fatalf("Expected zero date, got %v", zero)
		}
		if local.Day() != now.Day() {
			t.Fatalf("Day: Expected %d, got %d", now.Day(), local.Day())
		}
		if local.Month() != now.Month() {
			t.Fatalf("Month: Expected %d, got %d", now.Month(), local.Month())
		}
		if local.Year() != now.Year() {
			t.Fatalf("Year: Expected %d, got %d", now.Year(), local.Year())
		}
		if local.Weekday() != now.Weekday() {
			t.Fatalf("Weekday: Expected %d, got %d", now.Weekday(), local.Weekday())
		}
		if local.YearDay() != now.YearDay() {
			t.Fatalf("YearDay: Expected %d, got %d", now.YearDay(), local.YearDay())
		}

		y, w := local.ISOWeek()
		y2, w2 := now.ISOWeek()
		if y != y2 || w != w2 {
			t.Fatalf("ISOWeek mismatch")
		}

		year, m, d := local.Date()
		year2, m2, d2 := now.Date()
		if year != year2 || m != m2 || d != d2 {
			t.Fatalf("Date mismatch")
		}
	})
}

func TestDate(t *testing.T) {
	runBindingTest(t, DateOf, time.DateOnly, true)
	runDateBasicTest(t, DateOf)

	t.Run("Test Construct", func(t *testing.T) {
		now := time.Now()
		local := DateOf(now)
		if !local.Equal(NewDate(now.Year(), now.Month(), now.Day())) {
			t.Fatalf("Expected constructed date give different result")
		}
	})

	t.Run("Test Date Operations", func(t *testing.T) {
		now := NewDate(2023, 10, 15)
		nextDay := NewDate(2023, 10, 16)
		before := NewDate(2019, time.January, 1)
		after := NewDate(2023, 12, 15)

		if !now.Equal(NewDate(2023, 10, 15)) || now.Equal(nextDay) {
			t.Fatalf("Equally failed")
		}

		if sub := nextDay.Sub(now); sub != 24*time.Hour {
			t.Fatalf("Sub failed: expected 24 hours, got %v", nextDay.Sub(now))
		}
		if add := now.Add(1 * time.Hour); !add.Equal(now) {
			t.Fatalf("Add failed")
		}
		if add := now.Add(24 * time.Hour); !add.Equal(nextDay) {
			t.Fatalf("Add failed")
		}

		if addDays := now.AddDays(1); !addDays.Equal(nextDay) {
			t.Fatalf("AddDays failed")
		}
		if addDate := now.AddDate(0, 0, 1); !addDate.Equal(nextDay) {
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

		if dt := now.AtClock(5, 15, 20); !dt.Equal(DateTimeFrom(now, NewClock(5, 15, 20))) {
			t.Fatalf("AtClock produce wrong DateTime")
		}
		if dt := now.AtTime(5, 15, 20, 0); !dt.Equal(DateTimeFrom(now, NewTime(5, 15, 20, 0))) {
			t.Fatalf("AtTime produce wrong DateTime")
		}
		if !now.AtStartOfDay().Equal(DateTimeFrom(now, NewClock(0, 0, 0))) {
			t.Fatalf("AtStartOfDay produce wrong DateTime")
		}

		if date := NewDate(2025, 1, 1).UTC(); !date.Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)) {
			t.Fatalf("UTC produce wrong Date")
		}
		if date := NewDate(2025, 1, 1).Local(); !date.Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)) {
			t.Fatalf("Local produce wrong Date")
		}

		tz, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			t.Fatal(err)
		}
		if date := NewDate(2025, 1, 1).In(tz); !date.Equal(time.Date(2025, 1, 1, 0, 0, 0, 0, tz)) {
			t.Fatalf("In produce wrong Date")
		}
	})

	t.Run("Test ParseDate", func(t *testing.T) {
		tests := []struct {
			name string
			in   string
			want Date
		}{
			{"date only", "2023-10-15", NewDate(2023, 10, 15)},
			{"space datetime", "2023-10-15 12:10:15", NewDate(2023, 10, 15)},
			{"T datetime", "2023-10-15T12:10:15", NewDate(2023, 10, 15)},
			{"RFC3339Nano", "2023-10-15T12:10:15.123456789Z", NewDate(2023, 10, 15)},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ParseDate(tt.in)
				if err != nil {
					t.Fatalf("ParseDate(%q) unexpected error: %v", tt.in, err)
				}
				if !got.Equal(tt.want) {
					t.Fatalf("ParseDate(%q): expected %v, got %v", tt.in, tt.want, got)
				}
			})
		}

		for _, in := range []string{"", "2023", "2023-13-01", "2023-10-15T"} {
			if _, err := ParseDate(in); err == nil {
				t.Fatalf("ParseDate(%q) expected error, got nil", in)
			}
		}
	})
}
