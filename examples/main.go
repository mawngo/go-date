package main

import (
	"fmt"
	"github.com/mawngo/go-date/v2"
	"time"
)

func main() {
	d := date.NewDate(2026, 11, 14)
	d.AddDays(1)

	fmt.Println(d) // 2026-11-14
	// Convert to time.Time, at the start of the day in UTC.
	fmt.Println(d.UTC()) // 2026-11-14 00:00:00 +0000 UTC
	// Convert to time.Time, at the start of the day in local time.
	fmt.Println(d.Local()) // 2026-11-14 00:00:00 +0700 +07

	// Now
	fmt.Println(date.DateOf(time.Now())) // 2026-11-23
	// Next Monday
	fmt.Println(d.AtWeekday(time.Monday, 1)) // 2026-11-16
	// Prev Monday
	fmt.Println(d.AtWeekday(time.Monday, -1)) // 2026-11-09
	// Prev Saturday
	fmt.Println(d.AtWeekday(time.Saturday, -1)) // 2026-11-07
	// 2 week ago Saturday
	// d is Saturday, so if we -7, it falls right on the previous Saturday => we need to -8
	fmt.Println(d.AtWeekday(time.Saturday, -8)) // 2026-11-31

	// Convert to date.DateTime
	dt := d.AtClockNano(12, 0, 0, 0)
	fmt.Println(dt) // 2026-11-14 12:00:00
	// Then convert to time.Time
	fmt.Println(dt.Local()) // 2026-11-14 12:00:00 +0700 +07

	// 24h time only - out-of-range components are clamped into range.
	t := date.NewTimeNano(25, 70, 70, 1234000)
	fmt.Println(t)                      // 23:59:59.001234
	fmt.Println(t.Format(time.Kitchen)) // 11:59PM
	// Change time
	fmt.Println(dt.AtClockNano(t.ClockNano())) // 2026-11-14 23:59:59.001234
	// Or with second precision only
	fmt.Println(dt.AtClock(t.Clock())) // 2026-11-14 23:59:59
}
