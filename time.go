package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const timeFormat = time.TimeOnly + ".999999999"

// localTime is a generic interface for time types.
// Contains common methods for time types.
type localTime[T any] interface {
	Add(d time.Duration) T
	Sub(other T) time.Duration

	IsZero() bool
	Before(other T) bool
	After(other T) bool
	Equal(other T) bool
	Compare(other T) int

	Hour() int
	Minute() int
	Second() int
	Nanosecond() int

	Clock() (hour, minute, second int)
	AtDate(year int, month time.Month, day int) DateTime
	String() string
}

var _ localTime[Time] = (*Time)(nil)
var _ binding = (*Time)(nil)

// Time represents only the 24 hours wall-clock time part of the time.
// Current implementation uses [time.Time] to store the time - which is not efficient, we may consider using
// a more efficient representation in the future.
type Time struct {
	t time.Time
}

// NewTime create [Time] from hour, minute, second, and nanosecond.
// Time overflow will be wrapped around.
func NewTime(hour, min, sec, nsec int) Time {
	t := time.Date(1970, time.January, 1, hour, min, sec, nsec, time.UTC)
	return Time{
		// Set the date again to prevent the hour overflow changing the date.
		t: time.Date(1970, time.January, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC),
	}
}

// NewClock create [Time] from hour, minute and second, without nanosecond.
// Time overflow will be wrapped around.
func NewClock(hour, min, sec int) Time {
	return NewTime(hour, min, sec, 0)
}

// TimeOf create [Time] from [time.Time].
// Time overflow will be wrapped around.
func TimeOf(t time.Time) Time {
	if t.IsZero() {
		return Time{}
	}
	return NewTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// ClockOf create [Time] from [time.Time], without nanosecond.
// Time overflow will be wrapped around.
func ClockOf(t time.Time) Time {
	if t.IsZero() {
		return Time{}
	}
	return NewTime(t.Hour(), t.Minute(), t.Second(), 0)
}

// ParseTime parses a string and returns the time value it represents.
// ParseTime accepts an extended form of the RFC3339 partial-time format. After
// the HH:MM:SS part of the string, an optional fractional part may appear,
// consisting of a decimal point followed by one to nine decimal digits.
// (RFC3339 admits only one digit after the decimal point).
func ParseTime(s string) (Time, error) {
	t, err := time.Parse(timeFormat, s)
	if err != nil {
		return Time{}, err
	}
	return TimeOf(t), nil
}

// Sub returns the duration between [Time] and the other [Time].
// See [time.Time.Sub].
func (clock Time) Sub(other Time) time.Duration {
	t := clock.t.Sub(other.t)
	return t
}

// Add returns the [Time] after adding the duration d.
// The calculation is floored to day.
// See [time.Time.Add].
func (clock Time) Add(dur time.Duration) Time {
	t := clock.t.Add(dur)
	return NewTime(t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// IsZero reports whether [Time] represents the zero IsZero instants.
func (clock Time) IsZero() bool {
	return clock.t.IsZero()
}

// Before reports whether the [Time] is before others [Time].
func (clock Time) Before(other Time) bool {
	return clock.t.Before(other.t)
}

// After reports whether the [Time] is after others [Time].
func (clock Time) After(other Time) bool {
	return clock.t.After(other.t)
}

// Equal reports whether the [Time] is equal to others [Time].
func (clock Time) Equal(other Time) bool {
	return clock.t.Equal(other.t)
}

// Compare compares the time with u. If t is before u, it returns -1;
// if t is after u, it returns +1; if they're the same, it returns 0.
func (clock Time) Compare(other Time) int {
	return clock.t.Compare(other.t)
}

// Hour returns the hour component of the time, in the range [0, 23].
func (clock Time) Hour() int {
	return clock.t.Hour()
}

// Minute returns the minute component of the time, in the range [0, 59].
func (clock Time) Minute() int {
	return clock.t.Minute()
}

// Second returns the second component of the time, in the range [0, 59].
func (clock Time) Second() int {
	return clock.t.Second()
}

// Nanosecond returns the nanosecond component of the time.
func (clock Time) Nanosecond() int {
	return clock.t.Nanosecond()
}

// Clock returns the time components of clock as hour, minute, second.
func (clock Time) Clock() (hour, minute, second int) {
	return clock.t.Clock()
}

// Time returns the time components of clock as hour, minute, second, and nanosecond.
func (clock Time) Time() (hour, minute, second, nsec int) {
	hour, minute, second = clock.t.Clock()
	return hour, minute, second, clock.Nanosecond()
}

// AtDate returns a [DateTime] with the same time and the given date.
func (clock Time) AtDate(year int, month time.Month, day int) DateTime {
	return NewDateTime(year, month, day, clock.Hour(), clock.Minute(), clock.Second(), clock.Nanosecond())
}

func (clock Time) MarshalJSON() ([]byte, error) {
	if clock.IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(timeFormat)+len(`""`))
	b = append(b, '"')
	b = clock.t.AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (clock *Time) UnmarshalJSON(b []byte) error {
	data := string(b)
	if data == "null" {
		return nil
	}
	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("Time.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	t, err := ParseTime(data)
	if err != nil {
		return err
	}
	*clock = t
	return nil
}

func (clock Time) MarshalText() (text []byte, err error) {
	if clock.IsZero() {
		return nil, nil
	}
	return []byte(clock.t.Format(timeFormat)), nil
}

func (clock *Time) UnmarshalText(text []byte) error {
	param := string(text)
	if param == "" {
		return nil
	}
	t, err := ParseTime(param)
	if err != nil {
		return err
	}
	*clock = t
	return nil
}

// UnmarshalParam add support for gin param binding.
func (clock *Time) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	t, err := ParseTime(param)
	if err != nil {
		return err
	}
	*clock = t
	return nil
}

func (clock Time) GobEncode() ([]byte, error) {
	return clock.t.GobEncode()
}

func (clock *Time) GobDecode(b []byte) error {
	t := time.Time{}
	err := (&t).GobDecode(b)
	if err != nil {
		return err
	}
	*clock = TimeOf(t)
	return nil
}

func (clock Time) MarshalBinary() ([]byte, error) {
	return clock.t.MarshalBinary()
}

func (clock *Time) UnmarshalBinary(b []byte) error {
	t := time.Time{}
	err := (&t).UnmarshalBinary(b)
	if err != nil {
		return err
	}
	*clock = TimeOf(t)
	return nil
}

func (clock Time) Value() (driver.Value, error) {
	if clock.IsZero() {
		return nil, nil
	}
	return clock.t.Format(timeFormat), nil
}

func (clock *Time) Scan(value any) error {
	switch src := value.(type) {
	case nil:
		return nil
	case string:
		t, err := ParseTime(src)
		if err != nil {
			return err
		}
		*clock = t
	case []byte:
		t, err := ParseTime(string(src))
		if err != nil {
			return err
		}
		*clock = t
	case time.Time:
		*clock = TimeOf(src)
	default:
		return fmt.Errorf("Time.Scan: unsupported type: %T", src)
	}
	return nil
}

func (clock Time) String() string {
	return clock.t.Format(timeFormat)
}

func (clock Time) unwrap() time.Time {
	return clock.t
}

func (clock Time) GormDataType() string {
	return "time"
}
