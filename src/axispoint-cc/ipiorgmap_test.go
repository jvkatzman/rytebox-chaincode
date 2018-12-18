package main

import (
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var ipiOrg_in = `{"ipi":"JayZ","org":"Org1"}`
var ipiOrg_out = `{"docType":"IPIORGMAP","ipi":"JayZ","org":"Org1"}`

func Test_addIpiOrg(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addIpiOrg"), []byte(ipiOrg_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var ipiOrgKey = "JayZ"
	checkState(t, stub, ipiOrgKey, ipiOrg_out)

}
