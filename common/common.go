package common


// #OpType
const (
	READ = iota
	WRITE
)

type Operation struct {
	ClientId int
	Op int // #OpType
	Key string
	Value string
}

type OpTrace = []Operation
