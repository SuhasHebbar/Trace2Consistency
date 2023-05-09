package verifier

import (
	"cchkr/common"
	"cchkr/generator"
	"fmt"
	"testing"
)

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func TestNotSequential(t *testing.T) {
	// Client 1
	w13 := common.Operation{
		ClientId:   1,
		SequenceNo: 0,
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "3",
	}
	r12 := common.Operation{
		ClientId:   1,
		SequenceNo: 1,
		Op:         common.READ,
		Key:        "Key",
		Value:      "2",
	}
	r13 := common.Operation{
		ClientId:   1,
		SequenceNo: 2,
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "3",
	}
	c1 := common.OpTrace{
		w13,
		r12,
		r13,
	}

	// Client 2
	w22 := common.Operation{
		ClientId:   2,
		SequenceNo: 0,
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "2",
	}
	r23 := common.Operation{
		ClientId:   2,
		SequenceNo: 1,
		Op:         common.READ,
		Key:        "Key",
		Value:      "3",
	}
	r22 := common.Operation{
		ClientId:   2,
		SequenceNo: 2,
		Op:         common.READ,
		Key:        "Key",
		Value:      "2",
	}
	c2 := common.OpTrace{
		w22,
		r23,
		r22,
	}

	distTrace := common.DistTrace{
		1: c1,
		2: c2,
	}
	verifierCh := make(chan common.OpTrace, 1000)
	resultch := make(chan common.VerifierResult)
	g := generator.NewGenerator(distTrace, verifierCh)
	go g.RunGenerator()

	v := NewVerifier(verifierCh, resultch)
	go v.RunVerifier()

	result := <-v.resultCh

	for _, consistency := range result.ConsistencyProvided {
		fmt.Println(consistency)
	}

	for _, consistencyTrace := range result.Trace {
		fmt.Println(consistencyTrace)
	}

	if contains(result.ConsistencyProvided, "sequential") {
		t.Fatalf("Trace satisfies sequential but it shouldn't")
	}

}
