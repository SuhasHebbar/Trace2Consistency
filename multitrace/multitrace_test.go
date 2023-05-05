package multitrace

import (
	"cchkr/common"
	"testing"
)

func TestGetTraceConsistencies(t *testing.T) {
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

	output := GetTraceConsistencies(distTrace)
	if output.Exists("sequential consistency") {
		t.Fatalf("Trace satisfies sequential consistency but it shouldn't")
	}
}
