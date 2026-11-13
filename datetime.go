package date

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json/v2"
	"errors"
	"fmt"
	"time"
)

type binding interface {
	json.Marshaler
	json.Unmarshaler
	encoding.TextMarshaler
	encoding.TextUnmarshaler
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	gob.GobDecoder
	gob.GobEncoder
	driver.Valuer
	sql.Scanner
	UnmarshalParam(string) error
	GormDataType() string
	unwrap() time.Time

	String() string
	Format(layout string) string
	AppendFormat(b []byte, layout string) []byte
}

// localDate is a generic interface for date types.
// Contains common methods for date types.
type localDate[T any] interface {
	AddDate(years int, months int, days int) T
	AddDays(days int) T
	DaysSince(other T) int
	Sub(other T) time.Duration
	Add(d time.Duration) T

	IsZero() bool
	Before(other T) bool
	After(other T) bool
	Equal(other T) bool
	Compare(other T) int

	Day() int
	Month() time.Month
	Year() int
	Weekday() time.Weekday

	Date() (year int, month time.Month, day int)
	ISOWeek() (year, week int)
	YearDay() int

	AtTime(hour, minute, second, nsec int) DateTime
	AtClock(hour, minute, second int) DateTime
	AtStartOfDay() DateTime

	In(loc *time.Location) time.Time
	UTC() time.Time
	Local() time.Time
	String() string
}

var _ localDate[Date] = (*Date)(nil)
var _ binding = (*Date)(nil)
var _ localDate[DateTime] = (*DateTime)(nil)
var _ binding = (*DateTime)(nil)

// AtTime returns a new [DateTime] with the same date as date but with the given time.
// Time overflow can change the date.
func (date Date) AtTime(hour, minute, sec, nsec int) DateTime {
	return NewDateTime(date.t.Year(), date.t.Month(), date.t.Day(), hour, minute, sec, nsec)
}

// AtClock returns a new [DateTime] with the same date as date but with the given time without nanosecond.
// Time overflow can change the date.
func (date Date) AtClock(hour, minute, sec int) DateTime {
	return NewDateTime(date.t.Year(), date.t.Month(), date.t.Day(), hour, minute, sec, 0)
}

// AtStartOfDay returns [DateTime] at the start of the day.
func (date Date) AtStartOfDay() DateTime {
	//nolint:staticcheck
	return DateTime{
		t: date.t,
	}
}

// ----------------------------------------------------------------------------
// DateTime type.
// ----------------------------------------------------------------------------

// DateTime represents a local date time without zone.
type DateTime struct {
	t time.Time
}

// NewDateTime create [DateTime].
// Time overflow can change the date.
func NewDateTime(year int, month time.Month, day int, hour int, min int, sec int, nsec int) DateTime {
	return DateTime{
		t: time.Date(year, month, day, hour, min, sec, nsec, time.UTC),
	}
}

// DateTimeOf create [DateTime] from [time.Time].
func DateTimeOf(t time.Time) DateTime {
	if t.IsZero() {
		return DateTime{}
	}
	return NewDateTime(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// DateTimeNow create [DateTime] from [time.Now].
func DateTimeNow() DateTime {
	return DateTimeOf(time.Now())
}

// ParseDateTime parses a string and returns the DateTime it represents.
// Support [DateTimeNano] format, [time.RFC3339Nano] as a fallback format.
func ParseDateTime(s string) (DateTime, error) {
	t, err := time.Parse(DateTimeNano, s)
	if err == nil {
		return DateTimeOf(t), nil
	}
	if len(s) > 10 && (s[10] == 'T' || s[10] == 't') {
		sep := string(s[10])
		if t, err := time.Parse("2006-01-02"+sep+"15:04:05.999999999", s); err == nil {
			return DateTimeOf(t), nil
		}
		if t, err := time.Parse("2006-01-02"+sep+"15:04:05.999999999Z07:00", s); err == nil {
			return DateTimeOf(t), nil
		}
	}
	return DateTime{}, err
}

// Add returns the [DateTime] after adding the duration d.
// See [time.Time.Add].
func (dt DateTime) Add(d time.Duration) DateTime {
	t := dt.t.Add(d)
	return DateTimeOf(t)
}

// Sub returns the duration between [DateTime] and the other [DateTime].
// See [time.Time.Sub].
func (dt DateTime) Sub(other DateTime) time.Duration {
	return dt.t.Sub(other.t)
}

// Round returns the result of rounding t to the nearest multiple of d (since the zero time).
// The rounding behavior for halfway values is to round up.
// If d <= 0, Round returns t stripped of any monotonic clock reading but otherwise unchanged.
// See [time.Time.Round].
func (dt DateTime) Round(dur time.Duration) DateTime {
	return DateTimeOf(dt.t.Round(dur))
}

// Truncate returns the result of rounding t down to a multiple of d (since the zero time).
// If d <= 0, Truncate returns t stripped of any monotonic clock reading but otherwise unchanged.
// See [time.Time.Truncate].
func (dt DateTime) Truncate(dur time.Duration) DateTime {
	return DateTimeOf(dt.t.Truncate(dur))
}

// Before reports whether the [DateTime] is before others [DateTime].
func (dt DateTime) Before(other DateTime) bool {
	return dt.t.Before(other.t)
}

// After reports whether the [DateTime] is after others [DateTime].
func (dt DateTime) After(other DateTime) bool {
	return dt.t.After(other.t)
}

// Equal compare 2 datetime.
func (dt DateTime) Equal(other DateTime) bool {
	return dt.t.Equal(other.t)
}

// Compare compares the datetime with u. If t is before u, it returns -1;
// if t is after u, it returns +1; if they're the same, it returns 0.
func (dt DateTime) Compare(other DateTime) int {
	return dt.t.Compare(other.t)
}

// Hour returns the hour component of the time, in the range [0, 23].
func (dt DateTime) Hour() int {
	return dt.t.Hour()
}

// Minute returns the minute component of the time, in the range [0, 59].
func (dt DateTime) Minute() int {
	return dt.t.Minute()
}

// Second returns the second component of the time, in the range [0, 59].
func (dt DateTime) Second() int {
	return dt.t.Second()
}

// Nanosecond returns the nanosecond component of the datetime.
func (dt DateTime) Nanosecond() int {
	return dt.t.Nanosecond()
}

// Clock returns the time components of clock as hour, minute, second.
func (dt DateTime) Clock() (hour, minute, second int) {
	return dt.t.Clock()
}

// AtDate returns a [DateTime] with the same time and the given date.
func (dt DateTime) AtDate(year int, month time.Month, day int) DateTime {
	return NewDateTime(year, month, day, dt.Hour(), dt.Minute(), dt.Second(), dt.Nanosecond())
}

// AddDate add years, months, and days to date.
// See [time.Time.AddDate].
func (dt DateTime) AddDate(years int, months int, days int) DateTime {
	return DateTime{
		t: dt.t.AddDate(years, months, days),
	}
}

// AddDays add days to date.
func (dt DateTime) AddDays(days int) DateTime {
	return dt.AddDate(0, 0, days)
}

// DaysSince returns the signed number of days between the date and s, not including the end day.
// This is the inverse operation to AddDays.
func (dt DateTime) DaysSince(other DateTime) int {
	// Convert to start of day before calculating the difference - we only count the days difference.
	// After that, we convert to Unix time so we do not have to worry about leap seconds:
	// Unix time increases by exactly 86400 seconds per day.
	deltaUnix := dt.AtStartOfDay().t.Unix() - other.AtStartOfDay().t.Unix()
	return int(deltaUnix / 86400)
}

// IsZero reports whether [Date] represents the zero IsZero instants.
func (dt DateTime) IsZero() bool {
	return dt.t.IsZero()
}

// Day returns the day of the month specified by date.
func (dt DateTime) Day() int {
	return dt.t.Day()
}

// Month returns the month of the year specified by date.
func (dt DateTime) Month() time.Month {
	return dt.t.Month()
}

// Year returns the year in which date occurs.
func (dt DateTime) Year() int {
	return dt.t.Year()
}

// Weekday returns the day of the week specified by [Date].
// See [time.Time.Weekday].
func (dt DateTime) Weekday() time.Weekday {
	return dt.t.Weekday()
}

func (dt DateTime) Date() (year int, month time.Month, day int) {
	return dt.t.Date()
}

// ISOWeek returns the ISO 8601 year and week number in which t occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to
// week 52 or 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1
// of year n+1.
// See [time.Time.ISOWeek].
func (dt DateTime) ISOWeek() (year, week int) {
	return dt.t.ISOWeek()
}

// YearDay returns the day of the year specified by t, in the range [1,365] for non-leap years,
// and [1,366] in leap years.
// See [time.Time.YearDay].
func (dt DateTime) YearDay() int {
	return dt.t.YearDay()
}

// AtTime returns a new [DateTime] with the same date as date but with the given time.
// Time overflow can change the date.
func (dt DateTime) AtTime(hour, minute, second, nsec int) DateTime {
	return NewDateTime(dt.Year(), dt.Month(), dt.Day(), hour, minute, second, nsec)
}

// AtClock returns a new [DateTime] with the same date as date but with the given time without nanosecond.
// Time overflow can change the date.
func (dt DateTime) AtClock(hour, minute, second int) DateTime {
	return NewDateTime(dt.Year(), dt.Month(), dt.Day(), hour, minute, second, 0)
}

// AtStartOfDay returns a new [DateTime] at the start of day.
func (dt DateTime) AtStartOfDay() DateTime {
	return NewDateTime(dt.Year(), dt.Month(), dt.Day(), 0, 0, 0, 0)
}

func (dt DateTime) In(loc *time.Location) time.Time {
	return time.Date(dt.t.Year(), dt.t.Month(), dt.t.Day(), dt.t.Hour(), dt.t.Minute(), dt.t.Second(), dt.t.Nanosecond(), loc)
}

// UTC returns [time.Time] in UTC.
func (dt DateTime) UTC() time.Time {
	return dt.t
}

// Local returns [time.Time] in the local location.
func (dt DateTime) Local() time.Time {
	return time.Date(dt.t.Year(), dt.t.Month(), dt.t.Day(), dt.t.Hour(), dt.t.Minute(), dt.t.Second(), dt.t.Nanosecond(), time.Local)
}

func (dt DateTime) String() string {
	return dt.t.Format(DateTimeNano)
}

// Format returns a textual representation of the time value formatted according
// to the layout defined by the argument.
//
// See [time.Time.Format].
func (dt DateTime) Format(layout string) string {
	return dt.t.Format(layout)
}

// AppendFormat is like [Date.Format] but appends the textual
// representation to b and returns the extended buffer.
func (dt DateTime) AppendFormat(b []byte, layout string) []byte {
	return dt.t.AppendFormat(b, layout)
}

// ----------------------------------------------------------------------------
// Data bindings.
// ----------------------------------------------------------------------------

func (dt DateTime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(DateTimeNano)+len(`""`))
	b = append(b, '"')
	b = dt.t.AppendFormat(b, DateTimeNano)
	b = append(b, '"')
	return b, nil
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	data := string(b)
	if data == "null" {
		return nil
	}
	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("DateTime.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	t, err := ParseDateTime(data)
	if err != nil {
		return err
	}
	*dt = t
	return nil
}

func (dt DateTime) MarshalText() (text []byte, err error) {
	if dt.IsZero() {
		return nil, nil
	}
	return []byte(dt.t.Format(DateTimeNano)), nil
}

func (dt *DateTime) UnmarshalText(text []byte) error {
	param := string(text)
	if param == "" {
		return nil
	}
	t, err := ParseDateTime(param)
	if err != nil {
		return err
	}
	*dt = t
	return nil
}

func (dt DateTime) MarshalBinary() ([]byte, error) {
	return dt.t.MarshalBinary()
}

func (dt *DateTime) UnmarshalBinary(b []byte) error {
	t := time.Time{}
	err := (&t).UnmarshalBinary(b)
	if err != nil {
		return err
	}
	*dt = DateTimeOf(t)
	return nil
}

func (dt *DateTime) GobDecode(b []byte) error {
	t := time.Time{}
	err := (&t).GobDecode(b)
	if err != nil {
		return err
	}
	*dt = DateTimeOf(t)
	return nil
}

func (dt DateTime) GobEncode() ([]byte, error) {
	return dt.t.GobEncode()
}

func (dt DateTime) Value() (driver.Value, error) {
	if dt.IsZero() {
		return nil, nil
	}
	return dt.t.Format(DateTimeNano), nil
}

func (dt *DateTime) Scan(value any) error {
	switch src := value.(type) {
	case nil:
		return nil
	case string:
		t, err := ParseDateTime(src)
		if err != nil {
			return err
		}
		*dt = t
	case []byte:
		t, err := ParseDateTime(string(src))
		if err != nil {
			return err
		}
		*dt = t
	case time.Time:
		*dt = DateTimeOf(src)
	default:
		return fmt.Errorf("DateTime.Scan: unsupported type: %T", src)
	}
	return nil
}

func (dt *DateTime) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	t, err := ParseDateTime(param)
	if err != nil {
		return err
	}
	*dt = t
	return nil
}

func (dt DateTime) GormDataType() string {
	return "datetime"
}

func (dt DateTime) unwrap() time.Time {
	return dt.t
}
