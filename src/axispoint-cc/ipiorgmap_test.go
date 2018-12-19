package main

import (
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// *****************************************************************************
// ******************************* Mock Data ***********************************
// *****************************************************************************

var ipiOrg_in = `{"ipi":"JayZ","org":"org1"}`
var ipiOrg_out = `{"docType":"IPIORGMAP","ipi":"JayZ","org":"org1"}`

// *****************************************************************************

// *****************************************************************************
func MockIpiOrgResponse(functionName string) []byte {
	switch functionName {
	case "Test_AddIpiOrg_MappingExists":
		return []byte(`{"status":"500","message":"IPI-Org mapping already exists with this key: JayZ"}`)
	case "Test_UpdateIpiOrg_MappingExists":
		return []byte(`{"message": "IPI-Org mapping updated successfully"}`)
	case "Test_GetIpiOrgByUUID":
		return []byte(ipiOrg_out)
	default:
		return []byte("[]")
	}
}

// *****************************************************************************

func Test_AddIpiOrg(t *testing.T) {
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

func Test_AddIpiOrg_MappingExists(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addIpiOrg"), []byte(ipiOrg_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	//Calling Invoke a second time
	respPayload, err2 := checkInvoke(t, stub, [][]byte{[]byte("addIpiOrg"), []byte(ipiOrg_in)})
	if err2 != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var ipiOrgKey = "JayZ"
	checkState(t, stub, ipiOrgKey, ipiOrg_out)

	expected := MockIpiOrgResponse("Test_AddIpiOrg_MappingExists")
	if !reflect.DeepEqual(expected, respPayload) {
		t.Fatalf("Actual response is not equal to expected response")
	}

}

func Test_UpdateIpiOrg(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("updateIpiOrg"), []byte(ipiOrg_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var ipiOrgKey = "JayZ"
	checkState(t, stub, ipiOrgKey, ipiOrg_out)

}

func Test_UpdateIpiOrg_MappingExists(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("updateIpiOrg"), []byte(ipiOrg_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	//Calling Invoke a second time
	respPayload, err2 := checkInvoke(t, stub, [][]byte{[]byte("updateIpiOrg"), []byte(ipiOrg_in)})
	if err2 != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var ipiOrgKey = "JayZ"
	checkState(t, stub, ipiOrgKey, ipiOrg_out)

	expected := MockIpiOrgResponse("Test_UpdateIpiOrg_MappingExists")
	if !reflect.DeepEqual(expected, respPayload) {
		t.Fatalf("Actual response is not equal to expected response")
	}

}

func Test_GetIpiOrgByUUID(t *testing.T) {
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

	respPayload, err2 := testQuery(t, stub, "getIpiOrgByUUID", ipiOrgKey)
	if err2 != nil {
		t.Fatalf(err.Error())
	}

	expected := MockIpiOrgResponse("Test_GetIpiOrgByUUID")
	if !reflect.DeepEqual(expected, respPayload) {
		t.Fatalf("Actual response is not equal to expected response")
	}

}
