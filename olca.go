package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/msrocka/ld"
)

func toOlca(origin string, rels []*Rel) {

	processes := make(map[string]*ld.Process)
	flows := make(map[string]*ld.Flow)

	for _, rel := range rels {

		left := flows[rel.Left]
		if left == nil {
			left = makeFlow(rel.Left)
			flows[rel.Left] = left
		}

		proc := processes[rel.Left]
		if proc == nil {
			proc = makeProcess(rel.Left, left)
			processes[rel.Left] = proc
		}

		if rel.Left == rel.Right {
			// update the reference amount
			for i := range proc.Exchanges {
				if proc.Exchanges[i].QuantitativeReference {
					proc.Exchanges[i].Amount = rel.Amount
					break
				}
			}
			continue
		}

		right := flows[rel.Right]
		if right == nil {
			right = makeFlow(rel.Right)
			flows[rel.Right] = right
		}

		proc.Exchanges = append(proc.Exchanges, ld.Exchange{
			Amount: rel.Amount,
			Input:  rel.IsInput,
			Flow:   right.AsRef(),
			Unit: ld.NewReference(
				"Unit", "20aadc24-a391-41cf-b340-3e4529f44bde", "kg"),
			FlowProperty: ld.NewReference(
				"FlowProperty", "93a60a56-a3c8-11da-a746-0800200b9a66", "Mass")})
	}

	outPath := filepath.Base(origin)
	outPath = outPath + "_olca_" + time.Now().Format(time.RFC822) + ".zip"
	outPath = strings.ReplaceAll(strings.ReplaceAll(outPath, ":", "_"), " ", "_")
	fmt.Println("Write data to", outPath)

	w, err := ld.NewPackWriter(outPath)
	if err != nil {
		fmt.Println("ERROR: failed to create package", err)
		return
	}
	defer w.Close()
	for _, flow := range flows {
		err := w.PutFlow(flow)
		if err != nil {
			fmt.Println("ERROR: failed to write flow", flow.ID, err)
			return
		}
	}

	for _, process := range processes {
		err := w.PutProcess(process)
		if err != nil {
			fmt.Println("ERROR: failed to write process", process.ID, err)
			return
		}
	}
}

func makeFlow(id string) *ld.Flow {
	fmt.Println(" .. create flow", id)
	f := ld.Flow{}
	f.Name = id
	uid, err := uuid.NewRandom()
	if err != nil {
		f.ID = id
	} else {
		f.ID = uid.String()
	}

	if isWaste(id) {
		f.Type = ld.WasteFlow
	} else if isProduct(id) {
		f.Type = ld.ProductFlow
	} else {
		f.Type = ld.ElementaryFlow
	}

	f.FlowProperties = []ld.FlowPropertyFactor{
		ld.FlowPropertyFactor{
			ConversionFactor:      1.0,
			ReferenceFlowProperty: true,
			FlowProperty: ld.NewReference(
				"FlowProperty", "93a60a56-a3c8-11da-a746-0800200b9a66", "Mass")}}
	return &f
}

func makeProcess(id string, refFlow *ld.Flow) *ld.Process {
	fmt.Println(" .. create process", id)
	p := ld.Process{}
	p.Name = id
	uid, err := uuid.NewRandom()
	if err != nil {
		p.ID = id
	} else {
		p.ID = uid.String()
	}
	p.Type = ld.UnitProcess
	p.Exchanges = []ld.Exchange{
		ld.Exchange{
			Amount:                1.0,
			Input:                 refFlow.Type == ld.WasteFlow,
			QuantitativeReference: true,
			Flow:                  refFlow.AsRef(),
			Unit: ld.NewReference(
				"Unit", "20aadc24-a391-41cf-b340-3e4529f44bde", "kg"),
			FlowProperty: ld.NewReference(
				"FlowProperty", "93a60a56-a3c8-11da-a746-0800200b9a66", "Mass")}}
	return &p
}
