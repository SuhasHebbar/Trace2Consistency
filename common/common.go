package common

// #OpType
const (
	READ = iota
	WRITE
)

type Operation struct {
	ClientId   int
	SequenceNo int
	Op         int // #OpType
	Key        string
	Value      string
}

type OpTrace = []Operation
type DistTrace = map[int]OpTrace

type VerifierResult struct {
	ConsistencyProvided []string
	Trace               []OpTrace
}
