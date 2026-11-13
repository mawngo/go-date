# Go Date

Simple local `Date`, `Time` and `DateTime` type for Go.

This library provide very minimal functionality for date and datetime types, serving simple use cases.
For complex use cases, consider using [cloud.google.com/go/civil](https://pkg.go.dev/cloud.google.com/go/civil).

## Installation

Require go 1.27+

```shell
go get -u github.com/mawngo/go-date/v2
```

Or just copy the `date.go`, `datetime.go` and `time.go` file into your project.

## Usage

```go
package main

import (
	"fmt"
)

func main() {
	d := date.DateNow()
	d.AddDays(1)

	fmt.Println(d)
	// Convert to time.Time, at the start of the day in UTC.
	fmt.Println(d.UTC())
	// Convert to time.Time, at the start of the day in local time.
	fmt.Println(d.Local())
}

```

When unmarshalling JSON and text, `Date` and `DateTime` types support both their standard format and RFC3339Nano.