package multitrace

import (
	"bufio"
	"cchkr/common"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// See https://stackoverflow.com/a/16615559/12160191
func ParseFile(filename string) common.DistTrace {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintf("Failed to open file with name %v with error %v!", filename, err))
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	distTrace := common.DistTrace{}
	for scanner.Scan() {
		// parse each line
		ParseLine(scanner.Text(), distTrace)
	}

	return distTrace
}

// example line
// client=1 op=READ key=key value=value
func ParseLine(line string, distTrace map[int]common.OpTrace) {
	kv := ExtractKV(line)

	// parse the line into a struct

	// client=<client id>
	client, found := kv["client"]
	if !found {
		panic(fmt.Sprintf("client key not found in %v", line))
	}
	clientId, err := strconv.Atoi(client)
	if err != nil {
		panic(fmt.Sprintf("Invalid client id %v", client))
	}

	// determine the index in the program order
	seqno := len(distTrace[clientId])

	// op=READ or op=WRITE
	op, found := kv["op"]
	if !found {
		panic(fmt.Sprintf("op key not found in %v", line))
	}
	ope := common.READ
	if op == "WRITE" {
		ope = common.WRITE
	} else if op != "READ" {
		panic(fmt.Sprintf("Invalid op %v", op))
	}

	// key=<key>
	key, found := kv["key"]
	if !found {
		panic(fmt.Sprintf("key key not found in %v", line))
	}

	// value=<value>
	value, found := kv["value"]
	if !found {
		panic(fmt.Sprintf("value key not found in %v", line))
	}

	// Update dist trace by appending the following operation
	operation := common.Operation{
		ClientId:   clientId,
		SequenceNo: seqno,
		Op:         ope,
		Key:        key,
		Value:      value,
	}

	distTrace[clientId] = append(distTrace[clientId], operation)
}

// key value format
// key1=value1 key2=value2 ... keyn=valuen
func ExtractKV(line string) map[string]string {
	fields := strings.Fields(line)
	kv := map[string]string{}
	for _, field := range fields {
		// extract key, value from key=value
		kvfield := strings.Split(field, "=")
		if len(kvfield) != 2 {
			panic(fmt.Sprintf("Expected a key value pair but found %v", field))
		}

		key, value := kvfield[0], kvfield[1]
		_, found := kv[key]
		if found {
			panic(fmt.Sprintf("Duplicate key %v found", key))
		}

		kv[key] = value
	}

	return kv
}
