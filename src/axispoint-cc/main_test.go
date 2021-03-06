package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte, retval []byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
	if retval != nil {
		if res.Payload == nil {
			fmt.Printf("Init returned nil, expected %s", string(retval))
			t.FailNow()
		}
		if string(res.Payload) != string(retval) {
			fmt.Printf("Init returned %s, expected %s", string(res.Payload), string(retval))
			t.FailNow()
		}
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected, it was:", string(bytes))
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, args [][]byte, retval []byte) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Query", args[1], "failed", string(res.Message))
		t.FailNow()
	}
	if res.Payload == nil {
		fmt.Println("Query", args[1], "failed to get value")
		t.FailNow()
	}
	fmt.Println("Query Result", " was: ", string(res.Payload))
}

func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) ([]byte, error) {
	res := stub.MockInvoke("1", args)
	if res.Status != shim.OK {
		fmt.Println("Invoke", args, "failed", string(res.Message))
		return nil, errors.New(res.Message)
	}
	fmt.Println("Payload", string(res.Payload))
	return res.Payload, nil
}

func testQuery(t *testing.T, stub *shim.MockStub, function string, id string) ([]byte, error) {
	res := stub.MockInvoke("1", [][]byte{[]byte(function), []byte(id)})
	if res.Status != shim.OK {
		fmt.Println("Query", id, "failed", string(res.Message))
		return nil, errors.New(res.Message)
	}
	if res.Payload == nil {
		fmt.Println("Query", id, "failed to get value")
		return nil, errors.New(("Query failed to get value"))
	}
	fmt.Println("Query value", id, " is: ", string(res.Payload))
	return res.Payload, nil
}

func Test_Init(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("OK")}, nil)
}

func Test_ResetLedger(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("resetLedger"), []byte(exploitationReportSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	//verify # of records in the ledger are 0.
}
