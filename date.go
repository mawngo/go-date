package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const dateTimeFormat = time.DateTime + ".999999999"

// Date represents only date part of the time.
type Date struct {
	t time.Time
}

// NewDate create [Date] from year, month, and day.
func NewDate(year int, month time.Month, day int) Date {
	return Date{
		t: time.Date(year, month, day, 0, 0, 0, 0, time.UTC),
	}
}

// DateOf create [Date] from [time.Time].
func DateOf(t time.Time) Date {
	if t.IsZero() {
		return Date{}
	}
	return NewDate(t.Date())
}

// DateNow create [Date] from [time.Now].
func DateNow() Date {
	return DateOf(time.Now())
}

// ParseDate parses a string in full-date format and returns the date value it represents.
// Support [time.DateOnly] format and [time.DateTime], [YYYY-MM-DDTHH:MM:SS[.FFFFFFFFF]] and [time.RFC3339Nano] as fallback formats.
func ParseDate(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err == nil {
		return DateOf(t), nil
	}

	// Not very effective, but people should use the [time.DateOnly] instead.
	t, err2 := time.Parse(dateTimeFormat, s)
	if err2 == nil {
		return DateOf(t), nil
	}
	if len(s) > 10 && (s[10] == 'T' || s[10] == 't') {
		sep := string(s[10])
		if t, err := time.Parse("2006-01-02"+sep+"15:04:05.999999999", s); err == nil {
			return DateOf(t), nil
		}
		if t, err := time.Parse("2006-01-02"+sep+"15:04:05.999999999Z07:00", s); err == nil {
			return DateOf(t), nil
		}
	}
	return Date{}, err
}

// AddDate add years, months, and days to date.
// See [time.Time.AddDate].
func (date Date) AddDate(years int, months int, days int) Date {
	return Date{
		t: date.t.AddDate(years, months, days),
	}
}

// AddDays add days to date.
func (date Date) AddDays(days int) Date {
	return date.AddDate(0, 0, days)
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
// This is the inverse operation to AddDays.
func (date Date) DaysSince(other Date) (days int) {
	// We convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := date.t.Unix() - other.t.Unix()
	return int(deltaUnix / 86400)
}

// Sub returns the duration between [Date] and the other [Date].
// See [time.Time.Sub].
func (date Date) Sub(other Date) time.Duration {
	return date.t.Sub(other.t)
}

// Add returns the [Date] after adding the duration d.
// The calculation is floored to day.
// See [time.Time.Add].
func (date Date) Add(dur time.Duration) Date {
	t := date.t.Add(dur)
	return DateOf(t)
}

// Before reports whether the [Date] is before others [Date].
func (date Date) Before(other Date) bool {
	return date.t.Before(other.t)
}

// After reports whether the [Date] is after others [Date].
func (date Date) After(other Date) bool {
	return date.t.After(other.t)
}

// Equal compare 2 date.
func (date Date) Equal(other Date) bool {
	return date.t.Equal(other.t)
}

// Compare compares the date with u. If t is before u, it returns -1;
// if t is after u, it returns +1; if they're the same, it returns 0.
func (date Date) Compare(u Date) int {
	return date.t.Compare(u.t)
}

// UTC returns [time.Time] at the start of the day in UTC.
func (date Date) UTC() time.Time {
	return date.t
}

// In returns [time.Time] at the start of the day in the specified [time.Location].
func (date Date) In(loc *time.Location) time.Time {
	return time.Date(date.t.Year(), date.t.Month(), date.t.Day(), 0, 0, 0, 0, loc)
}

// Local returns [time.Time] at the start of the day in the local location.
func (date Date) Local() time.Time {
	return time.Date(date.t.Year(), date.t.Month(), date.t.Day(), 0, 0, 0, 0, time.Local)
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

// Weekday returns the day of the week specified by [Date].
// See [time.Time.Weekday].
func (date Date) Weekday() time.Weekday {
	return date.t.Weekday()
}

// ISOWeek returns the ISO 8601 year and week number in which t occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
// See [time.Time.ISOWeek].
func (date Date) ISOWeek() (year, week int) {
	return date.t.ISOWeek()
}

// YearDay returns the day of the year specified by t, in the range [1,365] for non-leap years,
// and [1,366] in leap years.
// See [time.Time.YearDay].
func (date Date) YearDay() int {
	return date.t.YearDay()
}

func (date Date) String() string {
	return date.t.Format(time.DateOnly)
}

func (date *Date) Scan(value any) error {
	switch src := value.(type) {
	case nil:
		return nil
	case string:
		t, err := ParseDate(src)
		if err != nil {
			return err
		}
		*date = t
	case []byte:
		t, err := ParseDate(string(src))
		if err != nil {
			return err
		}
		*date = t
	case time.Time:
		*date = DateOf(src)
	default:
		return fmt.Errorf("Date.Scan: unsupported type: %T", src)
	}
	return nil
}

func (date Date) Value() (driver.Value, error) {
	if date.IsZero() {
		return nil, nil
	}
	return date.t.Format(time.DateOnly), nil
}

func (date Date) GobEncode() ([]byte, error) {
	return date.t.GobEncode()
}

func (date *Date) GobDecode(b []byte) error {
	t := time.Time{}
	err := (&t).GobDecode(b)
	if err != nil {
		return err
	}
	*date = DateOf(t)
	return nil
}

func (date Date) MarshalBinary() ([]byte, error) {
	return date.t.MarshalBinary()
}

func (date *Date) UnmarshalBinary(b []byte) error {
	t := time.Time{}
	err := (&t).UnmarshalBinary(b)
	if err != nil {
		return err
	}
	*date = DateOf(t)
	return nil
}

func (date Date) MarshalJSON() ([]byte, error) {
	if date.IsZero() {
		return []byte("null"), nil
	}
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
		return errors.New("Date.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	t, err := ParseDate(data)
	if err != nil {
		return err
	}
	*date = t
	return nil
}

// UnmarshalParam add support for gin param binding.
func (date *Date) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	t, err := ParseDate(param)
	if err != nil {
		return err
	}
	*date = t
	return nil
}

func (date *Date) UnmarshalText(text []byte) error {
	param := string(text)
	if param == "" {
		return nil
	}
	t, err := ParseDate(param)
	if err != nil {
		return err
	}
	*date = t
	return nil
}

func (date Date) MarshalText() ([]byte, error) {
	if date.IsZero() {
		return nil, nil
	}
	return []byte(date.t.Format(time.DateOnly)), nil
}

func (date Date) GormDataType() string {
	return "date"
}

func (date Date) unwrap() time.Time {
	return date.t
}
