# Go Date

Simple local `Date`, `Time` and `DateTime` type for Go.

## Installation

Require go 1.27+

```shell
go get -u github.com/mawngo/go-date/v2
```

Alternatively, you can just copy the `date.go`, `datetime.go`, `time.go`, and `partial.go` files into your project.

The `datetime.go`, `time.go`, and `partial.go` files are optional; you can keep only the ones that you need.
`datetime.go` and `partial.go` depend on `date.go`, while `time.go` depends on `datetime.go`.

## Usage

```go
package main

import (
	"fmt"
)

func main() {
	d := date.NewDate(2026, 11, 14)
	d.AddDays(1)

	fmt.Println(d) // 2026-11-14
	// Convert to time.Time, at the start of the day in UTC.
	fmt.Println(d.UTC()) // 2026-11-14 00:00:00 +0000 UTC
	// Convert to time.Time, at the start of the day in local time.
	fmt.Println(d.Local()) // 2026-11-14 00:00:00 +0700 +07

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

```

See [examples/main.go](examples/main.go) for the full example.

When unmarshalling JSON and text, the `Date` and `DateTime` types support both their respective standard formats and
RFC3339Nano.