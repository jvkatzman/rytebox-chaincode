package main

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ********************************* Mock Data *********************************
var holderRepresentationSingle_in = `[{"holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}]`
var holderRepresentationSingle_out = `{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}`

var holderRepresentationMultiple_in = `[{"holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]},{"holderRepresentationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}]`
var holderRepresentationMultiple_out1 = `{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}`
var holderRepresentationMultiple_out2 = `{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}`

// *****************************************************************************
func MockGetHolderRepresentationResponse(functionName string) []byte {
	switch functionName {
	case "Test_AddHolderRepresentations_Single":
		return []byte(`{"successCount":1,"failureCount":0,"holderRepresentations":[]}`)
	case "Test_AddHolderRepresentations_Single_AlreadyExists":
		return []byte(`{"successCount":0,"failureCount":1,"holderRepresentations":[{"holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","message":"Holder Representation already exists!","success":false}]}`)
	case "Test_AddHolderRepresentations_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"holderRepresentations":[]}`)
	case "Test_GetHolderRepresentations":
		return []byte(`[{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]},{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}]`)
	case "Test_GetHolderRepresentationByUUID":
		return []byte(`{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}`)
	case "Test_GetHolderRepresentationByUUID_Failure":
		return []byte(`{"status":"500","message":"Holder Representation with UUID: 85fff2bf-00a2-423b-9567-55c6f4ee6ee2 does not exist"}`)
	default:
		return []byte("[]")
	}
}

func MockGetHolderRepresentationQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	if strings.Contains(queryString, "HOLDERREPRESENTATION") {
		return []string{`{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"85fff2bf-00a2-423b-9567-55c6f4ee6ee1","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]},{"docType":"HOLDERREPRESENTATION","holderRepresentationUUID":"817903a5-8a5f-4d51-af47-5bae33fc15b3","holder":"GECKOS!!","holderName":"KIERAN CASH","startDate":"20180131","endDate":"20181231","representations":[{"selector":"12.366","representative":"111111","representativeName":"test1"},{"selector":"12.366","representative":"222222","representativeName":"test2"}]}`}, nil
	}
	return nil, nil
}

func Test_AddHolderRepresentations_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var holderRepresentationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, holderRepresentationUUID, holderRepresentationSingle_out)

	expected := MockGetHolderRepresentationResponse("Test_AddHolderRepresentations_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddHolderRepresentations_Single_AlreadyExists(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Holder Representation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetHolderRepresentationResponse("Test_AddHolderRepresentations_Single_AlreadyExists")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddHolderRepresentations_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for first Exploitation Report
	var holderRepresentationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, holderRepresentationUUID, holderRepresentationMultiple_out1)

	// Check State for second Exploitation Report
	holderRepresentationUUID = "817903a5-8a5f-4d51-af47-5bae33fc15b3"
	checkState(t, stub, holderRepresentationUUID, holderRepresentationMultiple_out2)

	expected := MockGetHolderRepresentationResponse("Test_AddHolderRepresentations_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_AddHolderRepresentations_Empty1(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(`[]`)})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_AddHolderRepresentations_Empty2(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_GetHolderRepresentations(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	getHolderRepresentationsForQueryString = MockGetHolderRepresentationQueryResultForQueryString

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getHolderRepresentations"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetHolderRepresentationResponse("Test_GetHolderRepresentations")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetHolderRepresentationByUUID(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Holder Representation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getHolderRepresentationByUUID"), []byte("85fff2bf-00a2-423b-9567-55c6f4ee6ee1")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetHolderRepresentationResponse("Test_GetHolderRepresentationByUUID")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetHolderRepresentationByUUID_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Holder Representation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getHolderRepresentationByUUID"), []byte("85fff2bf-00a2-423b-9567-55c6f4ee6ee2")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetHolderRepresentationResponse("Test_GetHolderRepresentationByUUID_Failure")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateHolderRepresentations_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Holder Representation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateHolderRepresentations"), []byte(holderRepresentationSingle_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var holderRepresentationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, holderRepresentationUUID, holderRepresentationSingle_out)

	expected := MockGetHolderRepresentationResponse("Test_AddHolderRepresentations_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateHolderRepresentations_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	// Add Holder Representations
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addHolderRepresentations"), []byte(holderRepresentationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateHolderRepresentations"), []byte(holderRepresentationMultiple_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for first Exploitation Report
	var holderRepresentationUUID = "85fff2bf-00a2-423b-9567-55c6f4ee6ee1"
	checkState(t, stub, holderRepresentationUUID, holderRepresentationMultiple_out1)

	// Check State for second Exploitation Report
	holderRepresentationUUID = "817903a5-8a5f-4d51-af47-5bae33fc15b3"
	checkState(t, stub, holderRepresentationUUID, holderRepresentationMultiple_out2)

	expected := MockGetHolderRepresentationResponse("Test_AddHolderRepresentations_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}
