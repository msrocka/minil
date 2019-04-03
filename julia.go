package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func toJulia(origin string, rels []*Rel) {

	productIdx := make(map[string]int)
	products := make([]string, 0)
	eflowIdx := make(map[string]int)
	eflows := make([]string, 0)

	matrixA := make(map[int]map[int]float64)
	matrixB := make(map[int]map[int]float64)

	// index the products and elementary flows and
	// fill the matrix values
	for _, rel := range rels {

		if _, ok := productIdx[rel.Left]; !ok {
			productIdx[rel.Left] = len(productIdx)
			products = append(products, rel.Left)
		}

		col := productIdx[rel.Left]
		row := -1
		var matrix map[int]map[int]float64

		if isProduct(rel.Right) || isWaste(rel.Right) {
			if _, ok := productIdx[rel.Right]; !ok {
				productIdx[rel.Right] = len(productIdx)
				products = append(products, rel.Right)
			}
			row = productIdx[rel.Right]
			matrix = matrixA
		} else {
			if _, ok := eflowIdx[rel.Right]; !ok {
				eflowIdx[rel.Right] = len(eflowIdx)
				eflows = append(eflows, rel.Right)
			}
			row = eflowIdx[rel.Right]
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

	text := "# ref is the index of the reference flow\n"
	text += "# of the product system; 1 is the default value\n"
	text += "ref = 1\n\n"

	text += "# the technology matrix\n"
	text += "A = [\n"
	for row := 0; row < len(products); row++ {
		text += "    "
		for col := 0; col < len(products); col++ {
			val := 0.0
			if rowVals, ok := matrixA[row]; ok {
				val = rowVals[col]
			}
			if row == col && val == 0.0 {
				val = 1.0
			}
			t := fmt.Sprintf("%.2f  ", val)
			if val >= 0 {
				text += " "
			}
			text += t
		}
		text += ";  # " + products[row] + "\n"
	}
	text += "]\n"
	text += "println(\"\\n\\nA = \")\n"
	text += "display(A)\n\n"

	text += "# the intervention matrix\n"
	text += "B = [\n"
	for row := 0; row < len(eflows); row++ {
		text += "    "
		for col := 0; col < len(products); col++ {
			val := 0.0
			if rowVals, ok := matrixB[row]; ok {
				val = rowVals[col]
			}
			t := fmt.Sprintf("%.2f  ", val)
			if val >= 0 {
				text += " "
			}
			text += t
		}
		text += ";  # " + eflows[row] + "\n"
	}
	text += "]\n"
	text += "println(\"\\n\\nB = \")\n"
	text += "display(B)\n\n"

	text += "# the final demand\n"
	text += "f = zeros(size(A)[1])\n"
	text += "f[ref] = 1.0\n"
	text += "println(\"\\n\\nf = \")\n"
	text += "display(f)\n\n"

	text += "# the scaling vector\n"
	text += "s = A \\ f\n"
	text += "println(\"\\n\\ns = \")\n"
	text += "display(s)\n\n"

	text += "# the LCI result\n"
	text += "g = B * s\n"
	text += "println(\"\\n\\ng = \")\n"
	text += "display(g)\n\n"

	outPath := outputFile(origin, ".jl")
	fmt.Println("Write Julia model to", outPath)
	err := ioutil.WriteFile(outPath, []byte(text), os.ModePerm)
	if err != nil {
		fmt.Println("ERROR: failed to write file", err)
	}
}
