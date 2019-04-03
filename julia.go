package main

import "fmt"

func toJulia(origin string, rels []*Rel) {

	products := make(map[string]int)
	eflows := make(map[string]int)

	matrixA := make(map[int]map[int]float64)
	matrixB := make(map[int]map[int]float64)

	// index the products and elementary flows and
	// fill the matrix values
	for _, rel := range rels {

		if _, ok := products[rel.Left]; !ok {
			products[rel.Left] = len(products)
		}

		col := products[rel.Left]
		row := -1
		var matrix map[int]map[int]float64

		if isProduct(rel.Right) || isWaste(rel.Right) {
			if _, ok := products[rel.Right]; !ok {
				products[rel.Right] = len(products)
			}
			row = products[rel.Right]
			matrix = matrixA
		} else {
			if _, ok := eflows[rel.Right]; !ok {
				eflows[rel.Right] = len(eflows)
			}
			row = eflows[rel.Right]
			matrix = matrixB
		}

		rowVals := matrix[row]
		if rowVals == nil {
			rowVals = make(map[int]float64)
			matrix[row] = rowVals
		}
		amount := rel.Amount
		if rel.IsInput {
			amount = -amount
		}
		rowVals[col] = amount
	}

	text := "// the technology matrix\n"
	text += "A = [\n"

	fmt.Println(text)
}
