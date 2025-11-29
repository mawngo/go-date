# Go Date

Date type for Go.

## Installation

Require go 1.25+

```shell
go get -u github.com/mawngo/go-date
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/mawngo/go-date"
)

func main() {
	d := date.Now()
	d.AddDay(1)
	fmt.Println(d)
	// Convert to time.Time, at the start of the day.
	fmt.Println(d.ToLocalTime())
}
```