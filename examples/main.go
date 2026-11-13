package main

import (
	"fmt"
	"github.com/mawngo/go-date/v2"
)

func main() {
	d := date.DateNow()
	d.AddDays(1)

	fmt.Println(d) // 2026-11-14
	// Convert to time.Time, at the start of the day in UTC.
	fmt.Println(d.UTC()) // 2026-11-14 00:00:00 +0000 UTC
	// Convert to time.Time, at the start of the day in local time.
	fmt.Println(d.Local()) // 2026-11-14 00:00:00 +0700 +07

	// Convert to date.DateTime
	dt := d.AtTime(12, 0, 0, 0)
	fmt.Println(dt) // 2026-11-14 12:00:00
	// Then convert to time.Time
	fmt.Println(dt.Local()) // 2026-11-14 12:00:00 +0700 +07

	// 24h time only
	t := date.NewTime(25, 70, 70, 1234000)
	fmt.Println(t) // 02:11:10.001234
	// Change time
	fmt.Println(dt.AtTime(t.Time())) // 2026-11-14 02:11:10.001234
	// Or with second precision only
	fmt.Println(dt.AtClock(t.Clock())) // 2026-11-14 02:11:10
}
