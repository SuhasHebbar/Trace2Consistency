package verifier

import (
	"cchkr/common"
)

type Verifier struct {
	verifierCh chan common.OpTrace
	resultCh   chan common.VerifierResult
}

func NewVerifier(verifierCh chan common.OpTrace, resultCh chan common.VerifierResult) *Verifier {
	return &Verifier{verifierCh: verifierCh, resultCh: resultCh}
}

func (v *Verifier) RunVerifier() {
	numConsistencies := 0
	consistencies := map[string]bool{
		"sequential":                         false,
		"monotonic reads":                    false,
		"consistent prefix + read my writes": false,
		"read my writes":                     false,
		"eventual":                           false,
	}

	var verifierResult common.VerifierResult
	for trace := range v.verifierCh {
		if !checkBasicConsistency(trace) {
			continue
		}
		for consistency, done := range consistencies {
			if done {
				continue
			}
			switch consistency {
			case "sequential":
				if checkSequential(trace) {
					consistencies[consistency] = true
					numConsistencies++
					verifierResult.ConsistencyProvided = append(verifierResult.ConsistencyProvided, "sequential")
					verifierResult.Trace = append(verifierResult.Trace, trace)
				}
			case "monotonic reads":
				if checkMonotonicReads(trace) {
					consistencies[consistency] = true
					numConsistencies++
					verifierResult.ConsistencyProvided = append(verifierResult.ConsistencyProvided, "monotonic reads")
					verifierResult.Trace = append(verifierResult.Trace, trace)
				}
			case "consistent prefix + read my writes":
				if checkCPandReadMyWrites(trace) {
					consistencies[consistency] = true
					numConsistencies++
					verifierResult.ConsistencyProvided = append(verifierResult.ConsistencyProvided, "consistent prefix + read my writes")
					verifierResult.Trace = append(verifierResult.Trace, trace)
				}
			case "read my writes":
				if checkReadMyWrites(trace) {
					consistencies[consistency] = true
					numConsistencies++
					verifierResult.ConsistencyProvided = append(verifierResult.ConsistencyProvided, "read my writes")
					verifierResult.Trace = append(verifierResult.Trace, trace)
				}
			default:
				consistencies[consistency] = true
				numConsistencies++
				verifierResult.ConsistencyProvided = append(verifierResult.ConsistencyProvided, "eventual")
				verifierResult.Trace = append(verifierResult.Trace, trace)
			}

		}
		if numConsistencies == 5 {
			// for _, val := range consistencyProvided {
			// 	fmt.Println(val)
			// }
			v.resultCh <- verifierResult
			return
		}
	}
	v.resultCh <- verifierResult
	close(v.resultCh)
}

// Check if this is a valid trace (reads make sense in terms of previous writes)
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
func checkSequential(currTrace common.OpTrace) bool {
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

// Right now return true since CP is as good as eventual
func checkCPandReadMyWrites(currTrace common.OpTrace) bool {
	currSequenceNumber := make(map[int]map[string]int)
	for i := len(currTrace) - 1; i >= 0; i-- {
		opt := currTrace[i]
		if opt.Op == common.READ {
			_, isPresent := currSequenceNumber[opt.ClientId]
			if !isPresent {
				continue
			} else {
				sequenceNo, isPresent := currSequenceNumber[opt.ClientId][opt.Key]
				if !isPresent {
					continue
				}
				if opt.SequenceNo > sequenceNo {
					return false
				}
			}
		} else {
			_, isPresent := currSequenceNumber[opt.ClientId]
			if !isPresent {
				currSequenceNumber[opt.ClientId] = make(map[string]int)
				currSequenceNumber[opt.ClientId][opt.Key] = opt.SequenceNo
			} else {
				sequenceNo, isPresent := currSequenceNumber[opt.ClientId][opt.Key]
				if !isPresent || opt.SequenceNo < sequenceNo {
					currSequenceNumber[opt.ClientId][opt.Key] = opt.SequenceNo
				} else {
					return false
				}
			}
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
				if !isPresent || opt.SequenceNo < sequenceNo {
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
	for i := len(currTrace) - 1; i >= 0; i-- {
		opt := currTrace[i]
		if opt.Op == common.READ {
			_, isPresent := currSequenceNumber[opt.ClientId]
			if !isPresent {
				continue
			} else {
				sequenceNo, isPresent := currSequenceNumber[opt.ClientId][opt.Key]
				if !isPresent {
					continue
				}
				if opt.SequenceNo > sequenceNo {
					return false
				}
			}
		} else {
			_, isPresent := currSequenceNumber[opt.ClientId]
			if !isPresent {
				currSequenceNumber[opt.ClientId] = make(map[string]int)
				currSequenceNumber[opt.ClientId][opt.Key] = opt.SequenceNo

			} else {
				sequenceNo, isPresent := currSequenceNumber[opt.ClientId][opt.Key]
				if !isPresent || opt.SequenceNo < sequenceNo {
					currSequenceNumber[opt.ClientId][opt.Key] = opt.SequenceNo
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
