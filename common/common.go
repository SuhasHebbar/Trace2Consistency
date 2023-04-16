package common


// #OpType
const (
	READ = iota
	WRITE
)

type Operation struct {
	ClientId int
	SequenceNo int
	Op int // #OpType
	Key string
	Value string
}

type OpTrace = []Operation
