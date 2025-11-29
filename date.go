package date

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
	"time"
)

// Date represents only date part of the time.
type Date struct {
	t time.Time
}

// New create Date from year, month, and day.
func New(year int, month time.Month, day int) Date {
	return Date{
		t: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}

// FromTime convert time.Time to Date.
func FromTime(t time.Time) Date {
	y, m, d := t.Date()
	return Date{
		t: time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
	}
}

// Now create Date from time.Now.
func Now() Date {
	return FromTime(time.Now())
}

func (date Date) AddDate(years int, months int, days int) Date {
	return Date{
		t: date.t.AddDate(years, months, days),
	}
}

func (date Date) AddDay(days int) Date {
	return Date{
		t: date.t.AddDate(0, 0, days),
	}
}

// Before reports whether the Date is before others Date.
func (date Date) Before(other Date) bool {
	return date.t.Before(other.t)
}

// After reports whether the Date is after others Date.
func (date Date) After(other Date) bool {
	return date.t.After(other.t)
}

// Equal compare 2 date.
func (date Date) Equal(other Date) bool {
	return date.t.Equal(other.t)
}

// ToUTCTime convert Date to UTC time at the start of the day.
func (date Date) ToUTCTime() time.Time {
	return date.t
}

// ToLocationTime convert Date to time at the start of the day in the specified location.
func (date Date) ToLocationTime(loc *time.Location) time.Time {
	return time.Date(date.t.Year(), date.t.Month(), date.t.Day(), 0, 0, 0, 0, loc)
}

// ToLocalTime convert Date to time at the start of the day in the local location.
func (date Date) ToLocalTime() time.Time {
	return time.Date(date.t.Year(), date.t.Month(), date.t.Day(), 0, 0, 0, 0, time.Local)
}

// ToTime convert Date to time.Time.
func (date Date) ToTime(hour, minute, sec, nsec int, loc *time.Location) time.Time {
	return time.Date(date.t.Year(), date.t.Month(), date.t.Day(), hour, minute, sec, nsec, loc)
}

// IsZero reports whether [Date] represents the zero IsZero instants.
func (date Date) IsZero() bool {
	return date.t.IsZero()
}

// Day returns the day of the month specified by date.
func (date Date) Day() int {
	return date.t.Day()
}

// Month returns the month of the year specified by date.
func (date Date) Month() time.Month {
	return date.t.Month()
}

// Year returns the year in which date occurs.
func (date Date) Year() int {
	return date.t.Year()
}

// Date returns the year, month, and day in which date occurs.
func (date Date) Date() (year int, month time.Month, day int) {
	return date.t.Date()
}

func (date *Date) Scan(value interface{}) (err error) {
	nullTime := &sql.NullTime{}
	err = nullTime.Scan(value)
	*date = FromTime(nullTime.Time)
	return
}

func (date Date) Value() (driver.Value, error) {
	return date.t.Format(time.DateOnly), nil
}

func (date Date) GormDataType() string {
	return "date"
}

func (date Date) GobEncode() ([]byte, error) {
	return date.t.GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	return date.t.GobDecode(b)
}

func (date Date) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(time.DateOnly)+len(`""`))
	b = append(b, '"')
	b = date.t.AppendFormat(b, time.DateOnly)
	b = append(b, '"')
	return b, nil
}

func (date *Date) UnmarshalJSON(b []byte) error {
	data := string(b)
	if data == "null" {
		return nil
	}
	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("Time.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	d, err := time.Parse(time.DateOnly, data)
	if err != nil {
		if strings.Contains(data, "T") {
			err := date.t.UnmarshalJSON(b)
			if err != nil {
				return err
			}
			y, m, d := date.t.Date()
			date.t = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
			return nil
		}
		return err
	}
	date.t = d
	return nil
}

// UnmarshalParam add support for gin param binding.
func (date *Date) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	d, err := time.Parse(time.DateOnly, param)
	if err != nil {
		d, err = time.Parse(time.RFC3339, param)
		if err != nil {
			return err
		}
		y, m, day := d.Date()
		d = time.Date(y, m, day, 0, 0, 0, 0, time.UTC)
	}
	date.t = d
	return nil
}

func (date Date) String() string {
	return date.t.Format(time.DateOnly)
}
