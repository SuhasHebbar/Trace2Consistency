package main

import (
	"cchkr/common"
	"fmt"
)

func Concat(distTrace map[int32]common.OpTrace) common.OpTrace {
	serialTrace := common.OpTrace{}
	for _, trace := range distTrace {
		serialTrace = append(serialTrace, trace...)
	}
	return serialTrace
}

func checkBasicConsistency(currTrace common.OpTrace) bool {
	currWrittenVal := make(map[string]string)

	for _, opt := range currTrace {
		if opt.Op == common.READ {
			val, isPresent := currWrittenVal[opt.Key]
			if !isPresent || opt.Value != val {
				return false
			}
		} else {
			currWrittenVal[opt.Key] = opt.Value
		}
	}
	return true
}

// Reads and writes maintain program order
// Across Programs they can be interleaved
func checkSerializable(currTrace common.OpTrace) bool {
	currSequenceNumber := make(map[int]int)
	for _, opt := range currTrace {
		sequenceNo, isPresent := currSequenceNumber[opt.ClientId]
		if !isPresent || opt.SequenceNo > sequenceNo {
			currSequenceNumber[opt.ClientId] = opt.SequenceNo
		} else {
			return false
		}
	}
	return true
}

// Writes for a single object cannot be reordered
// Writes across objects can be interleaved
// Reads have no constraints
// All this reordering is w.r.t program order
func checkMonotonicReads(currTrace common.OpTrace) bool {
	currSequenceNumber := make(map[int]map[string]int)
	for _, opt := range currTrace {
		if opt.Op == common.WRITE {
			_, isPresent := currSequenceNumber[opt.ClientId]
			if !isPresent {
				currSequenceNumber[opt.ClientId] = make(map[string]int)
				currSequenceNumber[opt.ClientId][opt.Key] = opt.SequenceNo
			} else {
				sequenceNo, isPresent := currSequenceNumber[opt.ClientId][opt.Key]
				if !isPresent || opt.SequenceNo > sequenceNo {
					currSequenceNumber[opt.ClientId][opt.Key] = opt.SequenceNo
				} else {
					return false
				}
			}
		}
	}
	// for k, v := range currSequenceNumber {
	// 	fmt.Printf("ClientId : %v\n", k)
	// 	for k1, v1 := range v {
	// 		fmt.Printf("Key : %v\n", k1)
	// 		fmt.Printf("Sequence no. : %v\n", v1)
	// 	}
	// }
	return true
}

// Reads cannot be reordered before preceeding writes

func checkReadMyWrites(currTrace common.OpTrace) bool {
	currSequenceNumber := make(map[int]map[string]int)
	for _, opt := range currTrace {
		if opt.Op == common.READ {
			_, isPresent := currSequenceNumber[opt.ClientId]
			if !isPresent {
				return false
			} else {
				sequenceNo, isPresent := currSequenceNumber[opt.ClientId][opt.Key]
				if !isPresent || opt.SequenceNo <= sequenceNo {
					return false
				}
			}
		} else {
			_, isPresent := currSequenceNumber[opt.ClientId]
			if !isPresent {
				currSequenceNumber[opt.ClientId] = make(map[string]int)
			}
			currSequenceNumber[opt.ClientId][opt.Key] = opt.SequenceNo
		}
	}
	// for k, v := range currSequenceNumber {
	// 	fmt.Printf("ClientId : %v\n", k)
	// 	for k1, v1 := range v {
	// 		fmt.Printf("Key : %v\n", k1)
	// 		fmt.Printf("Sequence no. : %v\n", v1)
	// 	}
	// }
	return true
}

func verify() {
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
		Op:         common.WRITE,
		Key:        "Key",
		Value:      "2",
	}
	c2 := common.OpTrace{
		w22,
		r23,
		r22,
	}

	distTrace := map[int32]common.OpTrace{
		1: c1,
		2: c2,
	}

	currTrace := Concat(distTrace)

	// for i, opt := range serialTrace {
	// 	fmt.Printf("%v ...\n", i)
	// 	fmt.Println(opt.ClientId)
	// 	fmt.Println(opt.Key)
	// 	fmt.Println(opt.Value)
	// 	fmt.Println(opt.SequenceNo)
	// 	fmt.Println(opt.Op)
	// }

	fmt.Println(checkBasicConsistency(currTrace))

	fmt.Println(checkSerializable(currTrace))

	fmt.Println(checkMonotonicReads(currTrace))

	fmt.Println(checkReadMyWrites(currTrace))

}

func main() {
	verify()
}
