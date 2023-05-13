package main

import (
	"cchkr/common"
	"cchkr/generator"
	"cchkr/verifier"
	"fmt"
)

func contains(slice []string, value string) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

func main() {
	// Client 1
	w11 := common.Operation{
		ClientId:   1,
		SequenceNo: 0,
		Op:         common.WRITE,
		Key:        "K1",
		Value:      "1",
	}
	r12 := common.Operation{
		ClientId:   1,
		SequenceNo: 1,
		Op:         common.READ,
		Key:        "K1",
		Value:      "1",
	}
	w13 := common.Operation{
		ClientId:   1,
		SequenceNo: 2,
		Op:         common.WRITE,
		Key:        "K2",
		Value:      "3",
	}
	r14 := common.Operation{
		ClientId:   1,
		SequenceNo: 3,
		Op:         common.READ,
		Key:        "K2",
		Value:      "3",
	}
	c1 := common.OpTrace{
		w11,
		r12,
		w13,
		r14,
	}

	fmt.Println("Trace from Client 1 :")

	for _, val := range c1 {
		if val.Op == 0 {
			fmt.Printf("Read => ")
		} else {
			fmt.Printf("Write => ")
		}
		fmt.Printf("Key : %v, Val : %v\n", val.Key, val.Value)
	}

	// Client 2
	w21 := common.Operation{
		ClientId:   2,
		SequenceNo: 0,
		Op:         common.WRITE,
		Key:        "K1",
		Value:      "2",
	}
	r22 := common.Operation{
		ClientId:   2,
		SequenceNo: 1,
		Op:         common.READ,
		Key:        "K1",
		Value:      "2",
	}
	w23 := common.Operation{
		ClientId:   2,
		SequenceNo: 2,
		Op:         common.WRITE,
		Key:        "K2",
		Value:      "4",
	}
	r24 := common.Operation{
		ClientId:   2,
		SequenceNo: 3,
		Op:         common.READ,
		Key:        "K2",
		Value:      "4",
	}
	c2 := common.OpTrace{
		w21,
		r22,
		w23,
		r24,
	}

	fmt.Println("Trace from Client 2 :")

	for _, val := range c2 {
		if val.Op == 0 {
			fmt.Printf("Read => ")
		} else {
			fmt.Printf("Write => ")
		}
		fmt.Printf("Key : %v, Val : %v\n", val.Key, val.Value)
	}

	distTrace := map[int]common.OpTrace{
		1: c1,
		2: c2,
	}
	verifierCh := make(chan common.OpTrace, 1000)
	resultch := make(chan common.VerifierResult)

	fmt.Println("Running generator....")
	g := generator.NewGenerator(distTrace, verifierCh)
	go g.RunGenerator()

	fmt.Println("Running verifier....")
	v := verifier.NewVerifier(verifierCh, resultch)
	go v.RunVerifier()

	result := <-resultch

	idx := contains(result.ConsistencyProvided, "sequential")
	if idx != -1 {
		consistencyTrace := result.Trace[idx]

		fmt.Println("Given trace provides sequential consistency for the permutation :")

		for _, val := range consistencyTrace {
			if val.Op == 0 {
				fmt.Printf("Read => ")
			} else {
				fmt.Printf("Write => ")
			}
			fmt.Printf("Key : %v, Val : %v\n", val.Key, val.Value)
		}
	}
}
