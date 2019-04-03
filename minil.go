package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// Rel describes an input or output relation.
type Rel struct {
	Left    string
	Right   string
	Amount  float64
	IsInput bool
}

func isProduct(s string) bool {
	return strings.HasPrefix(s, "p") ||
		strings.HasPrefix(s, "P")
}

func isWaste(s string) bool {
	return strings.HasPrefix(s, "w") ||
		strings.HasPrefix(s, "W")
}

func parseLine(i int, line string) *Rel {
	fields := strings.Fields(line)
	if len(fields) != 4 {
		fmt.Println("Syntax error in line", i, ":", line)
		fmt.Println("  expected <ID> <DIR> <NUM> <ID>")
		return nil
	}

	left := fields[0]
	if !isProduct(left) && !isWaste(left) {
		fmt.Println("Syntax error in line", i, ":", line)
		fmt.Println("  left identifier is not a product or waste flow:", left)
		fmt.Println("  expected 'p' | 'P' | 'w' | 'W'")
		return nil
	}

	isInput := false
	switch fields[1] {
	case "->":
		isInput = false
	case "<-":
		isInput = true
	default:
		fmt.Println("Syntax error in line", i, ":", line)
		fmt.Println("  invalid direction ", fields[1])
		fmt.Println("  expected '<-' | '->'")
		return nil
	}

	amount, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		fmt.Println("Syntax error in line", i, ":", line)
		fmt.Println("  amount is not a valid number:", fields[2])
		return nil
	}

	right := fields[3]

	return &Rel{
		Left:    left,
		IsInput: isInput,
		Amount:  amount,
		Right:   right}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ERROR: no file given")
		return
	}

	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("ERROR: failed to open file", os.Args[1], err)
		return
	}

	text := string(bytes)
	var rels []*Rel
	for i, row := range strings.Split(text, "\n") {
		line := strings.TrimSpace(row)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		rel := parseLine(i, line)
		if rel == nil {
			// syntax error
			rels = nil
			break
		}
		rels = append(rels, rel)
	}

	if len(rels) == 0 {
		fmt.Println("Nothing to convert.")
		return
	}

	fmt.Println("Convert", len(rels), "relations ...")
	toOlca(os.Args[1], rels)

}
