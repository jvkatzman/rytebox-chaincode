package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ********************************* Mock Data *********************************
var ownerAdministrationSingle_in = `[{"ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}]`
var ownerAdministrationSingle_out = `{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}`

var ownerAdministrationMultiple_in = `[{"ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]},{"ownerAdministrationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}]`
var ownerAdministrationMultiple_out1 = `{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}`
var ownerAdministrationMultiple_out2 = `{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}`

// *****************************************************************************
func MockGetOwnerAdministrationResponse(functionName string) []byte {
	switch functionName {
	case "Test_AddOwnerAdministrations_Single":
		return []byte(`{"successCount":1,"failureCount":0,"ownerAdministrations":[]}`)
	case "Test_AddOwnerAdministrations_Single_AlreadyExists":
		return []byte(`{"successCount":0,"failureCount":1,"ownerAdministrations":[{"ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","message":"Owner Administration already exists!","success":false}]}`)
	case "Test_AddOwnerAdministrations_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"ownerAdministrations":[]}`)
	case "Test_GetOwnerAdministrations":
		return []byte(`[{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]},{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}]`)
	case "Test_GetOwnerAdministrationByUUID":
		return []byte(`{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}`)
	case "Test_GetOwnerAdministrationByUUID_Failure":
		return []byte(`{"status":"500","message":"UUID: 85fff2bf-00a2-423b-9567-55c6f4ee6ee2 does not exist"}`)
	default:
		return []byte("[]")
	}
}

func MockGetOwnerAdministrationQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	if strings.Contains(queryString, "OWNERADMINISTRATION") {
		return []string{`{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]},{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}`}, nil
	} else {
		return []string{`{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]},{"docType":"OWNERADMINISTRATION","ownerAdministrationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","owner":"GECKOS!!","ownerName":"KIERAN CASH","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-12-31T22:27:34.111Z","representations":[{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"111111","representativeName":"test1"},{"selector":"Territory in ('AUS', 'USA', 'GBR')","representative":"222222","representativeName":"test2"}]}`}, nil
	}
	return nil, nil
}

func Test_AddOwnerAdministrations_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var ownerAdministrationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, ownerAdministrationUUID, ownerAdministrationSingle_out)

	expected := MockGetOwnerAdministrationResponse("Test_AddOwnerAdministrations_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddOwnerAdministrations_Single_AlreadyExists(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Owner Administration
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetOwnerAdministrationResponse("Test_AddOwnerAdministrations_Single_AlreadyExists")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddOwnerAdministrations_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for first Exploitation Report
	var ownerAdministrationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, ownerAdministrationUUID, ownerAdministrationMultiple_out1)

	// Check State for second Exploitation Report
	ownerAdministrationUUID = "817903a5-8a5f-4d51-af47-5bae33fc15b3"
	checkState(t, stub, ownerAdministrationUUID, ownerAdministrationMultiple_out2)

	expected := MockGetOwnerAdministrationResponse("Test_AddOwnerAdministrations_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_AddOwnerAdministrations_Empty1(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(`[]`)})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_AddOwnerAdministrations_Empty2(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_GetOwnerAdministrations(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	getOwnerAdministrationsForQueryString = MockGetOwnerAdministrationQueryResultForQueryString

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getOwnerAdministrations"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetOwnerAdministrationResponse("Test_GetOwnerAdministrations")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetOwnerAdministrationByUUID(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Owner Administration
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getAssetByUUID"), []byte("85fff2bf-00a2-423b-9567-55c6f4ee6ee1")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetOwnerAdministrationResponse("Test_GetOwnerAdministrationByUUID")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetOwnerAdministrationByUUID_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Owner Administration
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getAssetByUUID"), []byte("85fff2bf-00a2-423b-9567-55c6f4ee6ee2")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetOwnerAdministrationResponse("Test_GetOwnerAdministrationByUUID_Failure")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateOwnerAdministrations_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Owner Administration
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateOwnerAdministrations"), []byte(ownerAdministrationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var ownerAdministrationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, ownerAdministrationUUID, ownerAdministrationSingle_out)

	expected := MockGetOwnerAdministrationResponse("Test_AddOwnerAdministrations_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateOwnerAdministrations_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	// Add Owner Administrations
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addOwnerAdministrations"), []byte(ownerAdministrationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateOwnerAdministrations"), []byte(ownerAdministrationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for first Exploitation Report
	var ownerAdministrationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, ownerAdministrationUUID, ownerAdministrationMultiple_out1)

	// Check State for second Exploitation Report
	ownerAdministrationUUID = "817903a5-8a5f-4d51-af47-5bae33fc15b3"
	checkState(t, stub, ownerAdministrationUUID, ownerAdministrationMultiple_out2)

	expected := MockGetOwnerAdministrationResponse("Test_AddOwnerAdministrations_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}
