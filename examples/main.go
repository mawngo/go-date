package main

import (
	"fmt"
	"github.com/mawngo/go-date"
	"time"
)

func main() {
	d := date.Now()
	d.AddDay(1)

	fmt.Println(d)
	fmt.Println(d.ToUTCTime())
	fmt.Println(d.ToLocalTime())
	fmt.Println(d.ToLocalTimeAtClock(time.Now().Clock()))
}
