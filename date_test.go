package date

import (
	"testing"
	"time"
)

func TestDate(t *testing.T) {
	t.Run("MarshalJson", func(t *testing.T) {
		b, err := FromTime(time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)).MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "\"2020-01-01\"" {
			t.Fatal("MarshalJson mismatch")
		}

		zone, err := time.LoadLocation("Asia/Ho_Chi_Minh")
		if err != nil {
			t.Fatal(err)
		}
		b, err = FromTime(time.Date(2020, 01, 01, 1, 1, 1, 1, zone)).MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "\"2020-01-01\"" {
			t.Fatal("MarshalJson mismatch")
		}
	})

	t.Run("UnmarshalJson", func(t *testing.T) {
		d := Date{}
		expect := FromTime(time.Date(2020, 01, 01, 1, 1, 1, 1, time.Local))

		err := (&d).UnmarshalJSON([]byte("\"2020-01-01\""))
		if err != nil {
			t.Fatal(err)
		}
		if d != expect {
			t.Fatal("UnmarshalJson mismatch")
		}

		err = (&d).UnmarshalJSON([]byte("\"2020-01-01T00:00:00+07:00\""))
		if err != nil {
			t.Fatal(err)
		}
		if d != expect {
			t.Fatal("UnmarshalParam mismatch")
		}
	})

	t.Run("GobEncode/Decode", func(t *testing.T) {
		expect := FromTime(time.Date(2020, 01, 01, 1, 1, 1, 1, time.Local))
		b, err := (&expect).GobEncode()
		if err != nil {
			t.Fatal(err)
		}

		actual := Date{}
		err = (&actual).GobDecode(b)
		if err != nil {
			t.Fatal(err)
		}
		if actual != expect {
			t.Fatal("Gob mismatch")
		}
	})

	t.Run("Marshal/UnmarshalText", func(t *testing.T) {
		expect := FromTime(time.Date(2020, 01, 01, 1, 1, 1, 1, time.Local))
		b, err := (&expect).MarshalText()
		if err != nil {
			t.Fatal(err)
		}
		if string(b) != "2020-01-01" {
			t.Fatal("MarshalText mismatch")
		}

		actual := Date{}
		err = (&actual).UnmarshalText(b)
		if err != nil {
			t.Fatal(err)
		}
		if actual != expect {
			t.Fatal("Gob mismatch")
		}
	})

	t.Run("UnmarshalParam", func(t *testing.T) {
		d := Date{}
		expect := FromTime(time.Date(2020, 01, 01, 1, 1, 1, 1, time.Local))

		err := (&d).UnmarshalParam("2020-01-01")
		if err != nil {
			t.Fatal(err)
		}
		if d != expect {
			t.Fatal("UnmarshalParam mismatch")
		}

		err = (&d).UnmarshalParam("2020-01-01T00:00:00+07:00")
		if err != nil {
			t.Fatal(err)
		}
		if d != expect {
			t.Fatal("UnmarshalParam mismatch")
		}
	})

	t.Run("Access Date", func(t *testing.T) {
		d := New(2020, 01, 01)
		if d.Year() != 2020 {
			t.Fatal("Year mismatch")
		}
		if d.Month() != time.Month(1) {
			t.Fatal("Month mismatch")
		}
		if d.Day() != 1 {
			t.Fatal("Day mismatch")
		}

		y, m, date := d.AddDay(-1).Date()
		if y != 2019 {
			t.Fatal("Year mismatch")
		}
		if m != time.Month(12) {
			t.Fatal("Month mismatch")
		}
		if date != 31 {
			t.Fatal("Day mismatch")
		}
	})
}
