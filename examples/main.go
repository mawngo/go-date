package main

import (
	"fmt"
	"github.com/mawngo/go-date"
)

func main() {
	d := date.DateNow()
	d.AddDays(1)

	fmt.Println(d)
	fmt.Println(d.UTC())
	fmt.Println(d.Local())
}
