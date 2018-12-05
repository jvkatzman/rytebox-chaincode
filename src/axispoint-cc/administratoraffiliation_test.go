package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ********************************* Mock Data *********************************
var administratorAffiliationSingle_in = `[{"administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}]`
var administratorAffiliationSingle_out = `{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}`

var administratorAffiliationMultiple_in = `[{"administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]},{"administratorAffiliationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}]`
var administratorAffiliationMultiple_out1 = `{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}`
var administratorAffiliationMultiple_out2 = `{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}`

// *****************************************************************************
func MockGetAdministratorAffiliationResponse(functionName string) []byte {
	switch functionName {
	case "Test_AddAdministratorAffiliations_Single":
		return []byte(`{"successCount":1,"failureCount":0,"administratorAffiliations":[]}`)
	case "Test_AddAdministratorAffiliations_Single_AlreadyExists":
		return []byte(`{"successCount":0,"failureCount":1,"administratorAffiliations":[{"administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","message":"Administrator Affiliation already exists!","success":false}]}`)
	case "Test_AddAdministratorAffiliations_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"administratorAffiliations":[]}`)
	case "Test_GetAdministratorAffiliations":
		return []byte(`[{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]},{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}]`)
	case "Test_GetAdministratorAffiliationByUUID":
		return []byte(`{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}`)
	case "Test_GetAdministratorAffiliationByUUID_Failure":
		return []byte(`{"status":"500","message":"Administrator Affiliation with UUID: 85fff2bf-00a2-423b-9567-55c6f4ee6ee2 does not exist"}`)
	default:
		return []byte("[]")
	}
}

func MockGetAdministratorAffiliationQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	if strings.Contains(queryString, "ADMINISTRATORAFFILIATION") {
		return []string{`{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]},{"docType":"ADMINISTRATORAFFILIATION","administratorAffiliationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","administrator":"GECKOS!!","startDate":"20180131","endDate":"20181231","affiliations":[{"selector":"12.366","affiliate":"111111","affiliateName":"test1"},{"selector":"12.366","affiliate":"222222","affiliateName":"test2"}]}`}, nil
	}
	return nil, nil
}

func Test_AddAdministratorAffiliations_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var administratorAffiliationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, administratorAffiliationUUID, administratorAffiliationSingle_out)

	expected := MockGetAdministratorAffiliationResponse("Test_AddAdministratorAffiliations_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddAdministratorAffiliations_Single_AlreadyExists(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Administrator Affiliation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetAdministratorAffiliationResponse("Test_AddAdministratorAffiliations_Single_AlreadyExists")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddAdministratorAffiliations_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for first Exploitation Report
	var administratorAffiliationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, administratorAffiliationUUID, administratorAffiliationMultiple_out1)

	// Check State for second Exploitation Report
	administratorAffiliationUUID = "817903a5-8a5f-4d51-af47-5bae33fc15b3"
	checkState(t, stub, administratorAffiliationUUID, administratorAffiliationMultiple_out2)

	expected := MockGetAdministratorAffiliationResponse("Test_AddAdministratorAffiliations_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_AddAdministratorAffiliations_Empty1(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(`[]`)})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_AddAdministratorAffiliations_Empty2(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_GetAdministratorAffiliations(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	getAdministratorAffiliationsForQueryString = MockGetAdministratorAffiliationQueryResultForQueryString

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getAdministratorAffiliations"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetAdministratorAffiliationResponse("Test_GetAdministratorAffiliations")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetAdministratorAffiliationByUUID(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Administrator Affiliation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getAdministratorAffiliationByUUID"), []byte("85fff2bf-00a2-423b-9567-55c6f4ee6ee1")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetAdministratorAffiliationResponse("Test_GetAdministratorAffiliationByUUID")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetAdministratorAffiliationByUUID_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Administrator Affiliation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getAdministratorAffiliationByUUID"), []byte("85fff2bf-00a2-423b-9567-55c6f4ee6ee2")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetAdministratorAffiliationResponse("Test_GetAdministratorAffiliationByUUID_Failure")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateAdministratorAffiliations_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Administrator Affiliation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateAdministratorAffiliations"), []byte(administratorAffiliationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var administratorAffiliationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, administratorAffiliationUUID, administratorAffiliationSingle_out)

	expected := MockGetAdministratorAffiliationResponse("Test_AddAdministratorAffiliations_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateAdministratorAffiliations_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	// Add Administrator Affiliations
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addAdministratorAffiliations"), []byte(administratorAffiliationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateAdministratorAffiliations"), []byte(administratorAffiliationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for first Exploitation Report
	var administratorAffiliationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, administratorAffiliationUUID, administratorAffiliationMultiple_out1)

	// Check State for second Exploitation Report
	administratorAffiliationUUID = "817903a5-8a5f-4d51-af47-5bae33fc15b3"
	checkState(t, stub, administratorAffiliationUUID, administratorAffiliationMultiple_out2)

	expected := MockGetAdministratorAffiliationResponse("Test_AddAdministratorAffiliations_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}
