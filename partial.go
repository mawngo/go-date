package date

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

const (
	// YearMonthOnly is the layout used to format and parse [YearMonth] in year-month format.
	YearMonthOnly = "2006-01"
	// MonthDayOnly is the layout used to format and parse [MonthDay] in month-day format.
	MonthDayOnly = "01-02"
)

var _ binding[YearMonth] = (*YearMonth)(nil)
var _ binding[MonthDay] = (*MonthDay)(nil)

// YearMonth represents only the year and month part of the time.
type YearMonth struct {
	t time.Time
}

// NewYearMonth creates a [YearMonth] from the year and month,
// clamping invalid values to the valid ranges defined by [ClampDate].
func NewYearMonth(year int, month time.Month) YearMonth {
	year, month, _ = ClampDate(year, month, 1)
	return YearMonth{
		t: time.Date(year, month, 1, 0, 0, 0, 0, time.UTC),
	}
}

// NewYearMonthChecked creates a [YearMonth] from the year and month,
// returning an error if any component is outside the valid ranges defined by [ClampDate].
func NewYearMonthChecked(year int, month time.Month) (YearMonth, error) {
	cyear, cmonth, _ := ClampDate(year, month, 1)
	if cyear != year {
		return YearMonth{}, fmt.Errorf("date.NewYearMonthChecked: invalid year %d (clamped into %d)", year, cyear)
	}
	if cmonth != month {
		return YearMonth{}, fmt.Errorf("date.NewYearMonthChecked: invalid month %d (clamped into %d)", month, cmonth)
	}
	return YearMonth{
		t: time.Date(year, month, 1, 0, 0, 0, 0, time.UTC),
	}, nil
}

// YearMonthOf create [YearMonth] from [time.Time].
func YearMonthOf(t time.Time) YearMonth {
	if t.IsZero() {
		return YearMonth{}
	}
	return NewYearMonth(t.Year(), t.Month())
}

// ParseYearMonth parses a string in [YearMonthOnly] format and returns the [YearMonth] it represents.
// Support [time.DateOnly] as a fallback format.
func ParseYearMonth(s string) (YearMonth, error) {
	t, err := time.Parse(YearMonthOnly, s)
	if err == nil {
		return YearMonthOf(t), nil
	}
	if t, err := time.Parse(time.DateOnly, s); err == nil {
		return YearMonthOf(t), nil
	}
	return YearMonth{}, err
}

// AddYearMonth add years and months to the [YearMonth].
// See [time.Time.AddDate].
func (ym YearMonth) AddYearMonth(years int, months int) YearMonth {
	date := ym.t.AddDate(years, months, 0)
	return NewYearMonth(date.Year(), date.Month())
}

// AddMonths add months to the YearMonth.
func (ym YearMonth) AddMonths(months int) YearMonth {
	return ym.AddYearMonth(0, months)
}

// MonthsSince returns the signed number of months between the [YearMonth] and other,
// This is the inverse operation to AddMonths.
func (ym YearMonth) MonthsSince(other YearMonth) (months int) {
	yearDiff := ym.t.Year() - other.t.Year()
	monthDiff := int(ym.t.Month()) - int(other.t.Month())
	return yearDiff*12 + monthDiff
}

// Sub returns the duration between [YearMonth] and the other [YearMonth].
// See [time.Time.Sub].
func (ym YearMonth) Sub(other YearMonth) time.Duration {
	return ym.t.Sub(other.t)
}

// Year returns the year in which [YearMonth] occurs.
func (ym YearMonth) Year() int {
	return ym.t.Year()
}

// Month returns the month of the year specified by [YearMonth].
func (ym YearMonth) Month() time.Month {
	return ym.t.Month()
}

// YearMonth returns the year and month in which [YearMonth] occurs.
func (ym YearMonth) YearMonth() (year int, month time.Month) {
	return ym.t.Year(), ym.t.Month()
}

// AtDay returns a new [Date] with the same year and month  as the [YearMonth] value on the given day,
// clamping the date components to the valid ranges defined by [ClampDate].
func (ym YearMonth) AtDay(day int) Date {
	return NewDate(ym.t.Year(), ym.t.Month(), day)
}

// IsValidDay returns true if the given day is valid for [YearMonth].
func (ym YearMonth) IsValidDay(day int) bool {
	y, m, d := ClampDate(ym.t.Year(), ym.t.Month(), day)
	return y == ym.t.Year() && m == ym.t.Month() && d == day
}

// AtStartOfMonth returns a [Date] with the same year and month as [YearMonth] on the first day of the month.
func (ym YearMonth) AtStartOfMonth() Date {
	return NewDate(ym.t.Year(), ym.t.Month(), 1)
}

// AtEndOfMonth returns a [Date] with the same year and month as [YearMonth] on the last day of the month.
func (ym YearMonth) AtEndOfMonth() Date {
	date := ym.t.AddDate(0, 1, -1)
	return NewDate(date.Year(), date.Month(), date.Day())
}

// String returns the year-month formatted in [YearMonthOnly].
func (ym YearMonth) String() string {
	return ym.t.Format(YearMonthOnly)
}

// Format returns a textual representation of the year-month value formatted according
// to the layout defined by the argument.
// The reference time is 00:00:00 UTC at the beginning of the month.
//
// See [time.Time.Format].
func (ym YearMonth) Format(layout string) string {
	return ym.t.Format(layout)
}

// AppendFormat is like [YearMonth.Format] but appends the textual
// representation to b and returns the extended buffer.
// The reference time is 00:00:00 UTC at the beginning of the month.
func (ym YearMonth) AppendFormat(b []byte, layout string) []byte {
	return ym.t.AppendFormat(b, layout)
}

// IsZero reports whether yearMonth represents the zero instant, January 1, year 1, 00:00:00 UTC.
// This matches [time.Time.IsZero]: a valid year-1 month is indistinguishable from the unset zero value.
func (ym YearMonth) IsZero() bool {
	return ym.t.IsZero()
}

// Before reports whether the [YearMonth] is before other [YearMonth].
func (ym YearMonth) Before(other YearMonth) bool {
	return ym.t.Before(other.t)
}

// After reports whether the [YearMonth] is after other [YearMonth].
func (ym YearMonth) After(other YearMonth) bool {
	return ym.t.After(other.t)
}

// Equal compare 2 year-month.
func (ym YearMonth) Equal(other YearMonth) bool {
	return ym.t.Equal(other.t)
}

// Compare compares the year-month with others. If ym is before other, it returns -1;
// if ym is after other, it returns +1; if they're the same, it returns 0.
func (ym YearMonth) Compare(other YearMonth) int {
	return ym.t.Compare(other.t)
}

// ----------------------------------------------------------------------------
// Data bindings.
// ----------------------------------------------------------------------------

func (ym YearMonth) MarshalJSON() ([]byte, error) {
	if ym.IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(YearMonthOnly)+len(`""`))
	b = append(b, '"')
	b = ym.t.AppendFormat(b, YearMonthOnly)
	b = append(b, '"')
	return b, nil
}

func (ym *YearMonth) UnmarshalJSON(b []byte) error {
	data := string(b)
	if data == "null" {
		return nil
	}
	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("YearMonth.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	t, err := ParseYearMonth(data)
	if err != nil {
		return err
	}
	*ym = t
	return nil
}

func (ym YearMonth) MarshalText() ([]byte, error) {
	if ym.IsZero() {
		return nil, nil
	}
	return []byte(ym.t.Format(YearMonthOnly)), nil
}

func (ym *YearMonth) UnmarshalText(text []byte) error {
	param := string(text)
	if param == "" {
		return nil
	}
	t, err := ParseYearMonth(param)
	if err != nil {
		return err
	}
	*ym = t
	return nil
}

func (ym YearMonth) MarshalBinary() ([]byte, error) {
	return ym.t.MarshalBinary()
}

func (ym *YearMonth) UnmarshalBinary(b []byte) error {
	t := time.Time{}
	err := (&t).UnmarshalBinary(b)
	if err != nil {
		return err
	}
	*ym = YearMonthOf(t)
	return nil
}

func (ym YearMonth) GobEncode() ([]byte, error) {
	return ym.t.GobEncode()
}

func (ym *YearMonth) GobDecode(b []byte) error {
	t := time.Time{}
	err := (&t).GobDecode(b)
	if err != nil {
		return err
	}
	*ym = YearMonthOf(t)
	return nil
}

func (ym YearMonth) Value() (driver.Value, error) {
	if ym.IsZero() {
		return nil, nil
	}
	return ym.t.Format(YearMonthOnly), nil
}

func (ym *YearMonth) Scan(value any) error {
	switch src := value.(type) {
	case nil:
		return nil
	case string:
		t, err := ParseYearMonth(src)
		if err != nil {
			return err
		}
		*ym = t
	case []byte:
		t, err := ParseYearMonth(string(src))
		if err != nil {
			return err
		}
		*ym = t
	case time.Time:
		*ym = YearMonthOf(src)
	default:
		return fmt.Errorf("YearMonth.Scan: unsupported type: %T", src)
	}
	return nil
}

// UnmarshalParam add support for gin param binding.
func (ym *YearMonth) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	t, err := ParseYearMonth(param)
	if err != nil {
		return err
	}
	*ym = t
	return nil
}

// ----------------------------------------------------------------------------
// Month Day.
// ----------------------------------------------------------------------------

// MonthDay represents only the month and day part of the time.
type MonthDay struct {
	t time.Time
}

// NewMonthDay creates a [MonthDay] from the month and day,
// clamping invalid values to the valid ranges defined by [ClampDate].
func NewMonthDay(month time.Month, day int) MonthDay {
	_, month, day = ClampDate(2000, month, day)
	return MonthDay{
		t: time.Date(2000, month, day, 0, 0, 0, 0, time.UTC),
	}
}

// NewMonthDayChecked creates a [MonthDay] from the month and day,
// returning an error if any component is outside the valid ranges defined by [ClampDate].
func NewMonthDayChecked(month time.Month, day int) (MonthDay, error) {
	_, cmonth, cday := ClampDate(2000, month, day)
	if cmonth != month {
		return MonthDay{}, fmt.Errorf("date.NewMonthDayChecked: invalid month %d (clamped into %d)", month, cmonth)
	}
	if cday != day {
		return MonthDay{}, fmt.Errorf("date.NewMonthDayChecked: invalid day %d (clamped into %d)", day, cday)
	}
	return MonthDay{
		t: time.Date(2000, month, day, 0, 0, 0, 0, time.UTC),
	}, nil
}

// MonthDayOf create [MonthDay] from [time.Time].
func MonthDayOf(t time.Time) MonthDay {
	if t.IsZero() {
		return MonthDay{}
	}
	return NewMonthDay(t.Month(), t.Day())
}

// ParseMonthDay parses a string in [MonthDayOnly] format and returns the [MonthDay] it represents.
// Support [time.DateOnly] as a fallback format.
func ParseMonthDay(s string) (MonthDay, error) {
	t, err := time.Parse(MonthDayOnly, s)
	if err == nil {
		return MonthDayOf(t), nil
	}
	if t, err := time.Parse(time.DateOnly, s); err == nil {
		return MonthDayOf(t), nil
	}
	return MonthDay{}, err
}

// Month returns the month of the year specified by [MonthDay].
func (md MonthDay) Month() time.Month {
	return md.t.Month()
}

// Day returns the day of the month specified by [MonthDay].
func (md MonthDay) Day() int {
	return md.t.Day()
}

// MonthDay returns the month and day in which [MonthDay] occurs.
func (md MonthDay) MonthDay() (month time.Month, day int) {
	return md.t.Month(), md.t.Day()
}

// AtYear returns a [Date] with the given year and the month and day from [MonthDay],
// clamping the date components to the valid ranges defined by [ClampDate].
//
// If the resulting date is invalid (for example, February 29 in a non-leap year),
// the clamped date is returned along with an error.
func (md MonthDay) AtYear(year int) Date {
	return NewDate(year, md.t.Month(), md.t.Day())
}

// IsValidYear returns true if the given year is valid for [MonthDay].
func (md MonthDay) IsValidYear(year int) bool {
	y, m, d := ClampDate(year, md.t.Month(), md.t.Day())
	return y == year && m == md.t.Month() && d == md.t.Day()
}

// String returns the month-day formatted in [MonthDayOnly].
func (md MonthDay) String() string {
	return md.t.Format(MonthDayOnly)
}

// Format returns a textual representation of the month-day value formatted according
// to the layout defined by the argument.
// The reference time is 00:00:00 UTC in the year 2000.
//
// See [time.Time.Format].
func (md MonthDay) Format(layout string) string {
	return md.t.Format(layout)
}

// AppendFormat is like [MonthDay.Format] but appends the textual
// representation to b and returns the extended buffer.
// The reference time is 00:00:00 UTC in the year 2000.
func (md MonthDay) AppendFormat(b []byte, layout string) []byte {
	return md.t.AppendFormat(b, layout)
}

// IsZero reports whether monthDay represents the zero instant.
func (md MonthDay) IsZero() bool {
	return md.t.IsZero()
}

// Before reports whether the [MonthDay] is before other [MonthDay].
func (md MonthDay) Before(other MonthDay) bool {
	return md.t.Before(other.t)
}

// After reports whether the [MonthDay] is after other [MonthDay].
func (md MonthDay) After(other MonthDay) bool {
	return md.t.After(other.t)
}

// Equal compare 2 month-day.
func (md MonthDay) Equal(other MonthDay) bool {
	return md.t.Equal(other.t)
}

// Compare compares the month-day with other. If md is before other, it returns -1;
// if md is after other, it returns +1; if they're the same, it returns 0.
func (md MonthDay) Compare(other MonthDay) int {
	return md.t.Compare(other.t)
}

// ----------------------------------------------------------------------------
// Data bindings.
// ----------------------------------------------------------------------------

func (md MonthDay) MarshalJSON() ([]byte, error) {
	if md.IsZero() {
		return []byte("null"), nil
	}
	b := make([]byte, 0, len(MonthDayOnly)+len(`""`))
	b = append(b, '"')
	b = md.t.AppendFormat(b, MonthDayOnly)
	b = append(b, '"')
	return b, nil
}

func (md *MonthDay) UnmarshalJSON(b []byte) error {
	data := string(b)
	if data == "null" {
		return nil
	}
	// TODO(https://go.dev/issue/47353): Properly unescape a JSON string.
	if len(data) < 2 || data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("MonthDay.UnmarshalJSON: input is not a JSON string")
	}
	data = data[len(`"`) : len(data)-len(`"`)]
	t, err := ParseMonthDay(data)
	if err != nil {
		return err
	}
	*md = t
	return nil
}

func (md MonthDay) MarshalText() ([]byte, error) {
	if md.IsZero() {
		return nil, nil
	}
	return []byte(md.t.Format(MonthDayOnly)), nil
}

func (md *MonthDay) UnmarshalText(text []byte) error {
	param := string(text)
	if param == "" {
		return nil
	}
	t, err := ParseMonthDay(param)
	if err != nil {
		return err
	}
	*md = t
	return nil
}

func (md MonthDay) MarshalBinary() ([]byte, error) {
	return md.t.MarshalBinary()
}

func (md *MonthDay) UnmarshalBinary(b []byte) error {
	t := time.Time{}
	err := (&t).UnmarshalBinary(b)
	if err != nil {
		return err
	}
	*md = MonthDayOf(t)
	return nil
}

func (md MonthDay) GobEncode() ([]byte, error) {
	return md.t.GobEncode()
}

func (md *MonthDay) GobDecode(b []byte) error {
	t := time.Time{}
	err := (&t).GobDecode(b)
	if err != nil {
		return err
	}
	*md = MonthDayOf(t)
	return nil
}

func (md MonthDay) Value() (driver.Value, error) {
	if md.IsZero() {
		return nil, nil
	}
	return md.t.Format(MonthDayOnly), nil
}

func (md *MonthDay) Scan(value any) error {
	switch src := value.(type) {
	case nil:
		return nil
	case string:
		t, err := ParseMonthDay(src)
		if err != nil {
			return err
		}
		*md = t
	case []byte:
		t, err := ParseMonthDay(string(src))
		if err != nil {
			return err
		}
		*md = t
	case time.Time:
		*md = MonthDayOf(src)
	default:
		return fmt.Errorf("MonthDay.Scan: unsupported type: %T", src)
	}
	return nil
}

// UnmarshalParam add support for gin param binding.
func (md *MonthDay) UnmarshalParam(param string) error {
	if param == "" {
		return nil
	}
	t, err := ParseMonthDay(param)
	if err != nil {
		return err
	}
	*md = t
	return nil
}
