package main

import (
	"cchkr/common"
	"cchkr/generator"
	"cchkr/verifier"
	"fmt"
	"strconv"
	"time"
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
	// Record start time
	startTime := time.Now()

	// // Run something
	// time.Sleep(500 * time.Millisecond)
	// Client 1
	var w1 []common.Operation
	for i := 0; i < 2; i++ {
		w1 = append(w1, common.Operation{
			ClientId:   1,
			SequenceNo: 0,
			Op:         common.WRITE,
			Key:        "K1",
			Value:      strconv.Itoa(i + 1),
		})
	}
	var r1 []common.Operation
	for i := 0; i < 2; i++ {
		r1 = append(r1, common.Operation{
			ClientId:   1,
			SequenceNo: 0,
			Op:         common.READ,
			Key:        "K1",
			Value:      strconv.Itoa(2 - i),
		})
	}
	c1 := common.OpTrace{}

	for i := 0; i < 2; i++ {
		c1 = append(c1, r1[i])
	}

	for i := 0; i < 2; i++ {
		c1 = append(c1, w1[i])
	}

	var w2 []common.Operation
	for i := 2; i < 4; i++ {
		w2 = append(w2, common.Operation{
			ClientId:   1,
			SequenceNo: 0,
			Op:         common.WRITE,
			Key:        "K1",
			Value:      strconv.Itoa(i + 1),
		})
	}

	var r2 []common.Operation
	for i := 2; i < 4; i++ {
		r2 = append(r2, common.Operation{
			ClientId:   1,
			SequenceNo: 0,
			Op:         common.READ,
			Key:        "K1",
			Value:      strconv.Itoa(4 - i),
		})
	}
	c2 := common.OpTrace{}

	for i := 0; i < 2; i++ {
		c2 = append(c2, r2[i])
	}

	for i := 0; i < 2; i++ {
		c2 = append(c2, w2[i])
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

	// Record end time
	endTime := time.Now()

	// Calculate elapsed time in milliseconds
	elapsedTime := endTime.Sub(startTime).Microseconds()

	fmt.Printf("Elapsed time: %d microseconds\n", elapsedTime)

	idx := contains(result.ConsistencyProvided, "eventual")
	if idx != -1 {
		//consistencyTrace := result.Trace[idx]

		fmt.Println("Given trace provides eventual consistency for the permutation :")

		// for _, val := range consistencyTrace {
		// 	if val.Op == 0 {
		// 		fmt.Printf("Read => ")
		// 	} else {
		// 		fmt.Printf("Write => ")
		// 	}
		// 	fmt.Printf("Key : %v, Val : %v\n", val.Key, val.Value)
		// }
	}
}
