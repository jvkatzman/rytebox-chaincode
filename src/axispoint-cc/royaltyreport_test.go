package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// ********************************* Mock Data *********************************
var exploitationReportSingle1_in = `[{"source":"M86321","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","isrc":"00029521","units":156062,"exploitationDate":"201811","amount":"36518.51","usageType":"SDIGM","exploitationReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eff","territory":"AUS"}]`
var exploitationReportSingle2_in = `[{"source":"M86321","songTitle":"LIVING WITH THE LAW","writerName":"CHRIS WHITLEY","isrc":"00029522","units":156062,"exploitationDate":"201811","amount":"36518.51","usageType":"SDIGM","exploitationReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eff","territory":"AUS"}]`

var royaltyReportSingle1_in = `[{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","source":"M86321","isrc":"00029521","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","units":156062,"exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}]`
var royaltyReportSingle1_out = `{"docType":"ROYALTYREPORT","royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","exploitationReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eff","source":"M86321","isrc":"00029521","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","units":156062,"exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}`
var royaltyReportMultiple1_in = `[{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","source":"M86321","isrc":"00029521","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","units":156062,"exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"},{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8ff","source":"M86321","isrc":"00029522","songTitle":"LIVING WITH THE LAW","writerName":"CHRIS WHITLEY","units":156062,"exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}]`
var royaltyReportSingle2_out = `{"docType":"ROYALTYREPORT","royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8ff","exploitationReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eff","source":"M86321","isrc":"00029522","songTitle":"LIVING WITH THE LAW","writerName":"CHRIS WHITLEY","units":156062,"exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}`
var royaltyReportSingle2_in = `[{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","source":"M86321","isrc":"00029523","exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}]`
var royaltyReportMultiple2_in = `[{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","source":"M86321","isrc":"00029524","exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"},{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8ff","source":"M86321","isrc":"00029525","exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}]`

// *****************************************************************************
func MockGetExploitationReportQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	if strings.Contains(queryString, "00029521") {
		return []string{`{"docType":"EXPLOITATIONREPORT","source":"M86321","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","isrc":"00029521","units":156062,"exploitationDate":"201811","amount":"36518.51","usageType":"SDIGM","exploitationReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eff","territory":"AUS"}`}, nil
	} else if strings.Contains(queryString, "00029522") {
		return []string{`{"docType":"EXPLOITATIONREPORT","source":"M86321","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","isrc":"00029522","units":156062,"exploitationDate":"201811","amount":"36518.51","usageType":"SDIGM","exploitationReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eff","territory":"AUS"}`}, nil
	}
	return nil, nil
}

func MockGetRoyaltyReportResponse(functionName string) []byte {
	switch functionName {
	case "Test_AddRoyaltyReports_Single":
		return []byte(`{"successCount":1,"failureCount":0,"royaltyReports":[]}`)
	case "Test_AddRoyaltyReports_Single_Failure":
		return []byte(`{"successCount":0,"failureCount":1,"royaltyReports":[{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","message":"Cannot find Exploitation Report with Source: M86321, ISRC: 00029523, Exploitation Date: 20170131, Territory: AUS, Usage Type: SDIGM","success":false}]}`)
	case "Test_AddRoyaltyReports_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"royaltyReports":[]}`)
	case "Test_AddRoyaltyReports_Multiple_Failure":
		return []byte(`{"successCount":0,"failureCount":2,"royaltyReports":[{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","message":"Cannot find Exploitation Report with Source: M86321, ISRC: 00029524, Exploitation Date: 20170131, Territory: AUS, Usage Type: SDIGM","success":false},{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8ff","message":"Cannot find Exploitation Report with Source: M86321, ISRC: 00029525, Exploitation Date: 20170131, Territory: AUS, Usage Type: SDIGM","success":false}]}`)
	default:
		return []byte("[]")
	}
}

func MockGetRoyaltyReportForPeriod(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	return []string{`{"docType":"ROYALTYREPORT","royaltyReportUUID":"2e0d3d18-f25f-36d0-81ef-6e3da893d4aa","exploitationReportUUID":"7cb134b2-156f-32e0-a4d9-6165c6ad1aca","source":"P8819H","isrc":"00055524","songTitle":"GECKOS!!","writerName":"KIERAN CASH","units":226,"exploitationDate":"20170131","amount":"18.5034375","rightType":"PERF","territory":"AUS","usageType":"SDIGP","target":"W998"},{"docType":"ROYALTYREPORT","royaltyReportUUID":"7f384cbf-0d0d-3698-9714-841b8ecb73f9","exploitationReportUUID":"6874280d-2897-3321-b238-0b4dfa0aa516","source":"P8819H","isrc":"00055524","songTitle":"GECKOS!!","writerName":"KIERAN CASH","units":150,"exploitationDate":"20170131","amount":"12.404250000000001","rightType":"MECH","territory":"AUS","usageType":"SDIGM","target":"W998"},{"docType":"ROYALTYREPORT","royaltyReportUUID":"94c878c5-f754-3b04-b90e-4f01cbd54ad6","exploitationReportUUID":"8ab33826-399f-3707-a0af-dfedc3d3b7f3","source":"P8819H","isrc":"00055524","songTitle":"GECKOS!!","writerName":"KIERAN CASH","units":164,"exploitationDate":"20170131","amount":"12.366","rightType":"MECH","territory":"AUS","usageType":"SMECH","target":"W998"}`}, nil
}

func Test_AddRoyaltyReports_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addExploitationReports"), []byte(exploitationReportSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	getExploitationReportForQueryString = MockGetExploitationReportQueryResultForQueryString
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyReports"), []byte(royaltyReportSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var royalReportUUID = "a4c7408b-d68b-499e-8dfa-ff81b43ca8fe"
	checkState(t, stub, royalReportUUID, royaltyReportSingle1_out)

	expected := MockGetRoyaltyReportResponse("Test_AddRoyaltyReports_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddRoyaltyReports_Single_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getExploitationReportForQueryString = MockGetExploitationReportQueryResultForQueryString
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyReports"), []byte(royaltyReportSingle2_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetRoyaltyReportResponse("Test_AddRoyaltyReports_Single_Failure")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddRoyaltyReports_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addExploitationReports"), []byte(exploitationReportSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}
	_, err = checkInvoke(t, stub, [][]byte{[]byte("addExploitationReports"), []byte(exploitationReportSingle2_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	getExploitationReportForQueryString = MockGetExploitationReportQueryResultForQueryString
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyReports"), []byte(royaltyReportMultiple1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var royalReportUUID = "a4c7408b-d68b-499e-8dfa-ff81b43ca8fe"
	checkState(t, stub, royalReportUUID, royaltyReportSingle1_out)

	royalReportUUID = "a4c7408b-d68b-499e-8dfa-ff81b43ca8ff"
	checkState(t, stub, royalReportUUID, royaltyReportSingle2_out)

	expected := MockGetRoyaltyReportResponse("Test_AddRoyaltyReports_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddRoyaltyReports_Multiple_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getExploitationReportForQueryString = MockGetExploitationReportQueryResultForQueryString
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyReports"), []byte(royaltyReportMultiple2_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetRoyaltyReportResponse("Test_AddRoyaltyReports_Multiple_Failure")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

// place holder -mock stub not implemented to read data from couch db
func Test_getRoyaltyDataForPeriod(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getRoyaltyReportsForQueryString = MockGetRoyaltyReportForPeriod
	actualOutput, err := checkInvoke(t, stub, [][]byte{[]byte("getRoyaltyDataForPeriod"), []byte("20170131"), []byte("M86322")})
	if err != nil {
		t.Fatalf(err.Error())
	}
	mockOutput, err := MockGetRoyaltyReportForPeriod(stub, "")
	//we need this becuase the data types are different -- normalize to string
	expectedOutput := fmt.Sprintf("[%s]", mockOutput[0])

	if !reflect.DeepEqual(expectedOutput, string(actualOutput)) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}
