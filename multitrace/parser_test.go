package multitrace

import (
	"cchkr/common"
	"reflect"
	"testing"
)

// See https://stackoverflow.com/a/18211675/12160191

func TestExtractKV(t *testing.T) {
	expected := map[string]string{
		"client": "1",
		"op":     "READ",
		"key":    "key",
		"value":  "value",
	}
	actual := ExtractKV("client=1 op=READ key=key value=value")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, Got %v", expected, actual)
	}
}

func TestParseLine(t *testing.T) {
	operation := common.Operation{
		ClientId:   1,
		SequenceNo: 0,
		Op:         common.READ,
		Key:        "key",
		Value:      "value",
	}
	expected := map[int]common.OpTrace{
		1: []common.Operation{operation},
	}
	distTrace := map[int]common.OpTrace{}
	actual := ParseLine("client=1 op=READ key=key value=value", distTrace)
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, Got %v", expected, actual)
	}
}
