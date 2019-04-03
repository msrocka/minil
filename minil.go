package main

import (
	"fmt"
	"os"
)

// Rel describes an input or output relation.
type Rel struct {
	Left    string
	Right   string
	Amount  float64
	IsInput bool
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no file given")
		return
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: failed to open file", os.Args[1], err)
		return
	}
	defer f.Close()

}
