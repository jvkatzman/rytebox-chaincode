package main

import (
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var collectionRightReportUUID = "15094dbb-9853-4737-aaa6-544ed27e0ac1"
var collectionRightReportSingleInput = `[{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"15094dbb-9853-4737-aaa6-544ed27e0ac1","from":"PU200004","fromName":"MARS FORCE MUSIC","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}]`
var collectionRightReportMultipleInput = `[{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"04240be9-73d3-4227-88a1-31c52d4db3bc","from":"PA300002","fromName":"URBAN SONGS","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]},{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"15094dbb-9853-4737-aaa6-544ed27e0ac1","from":"PU200004","fromName":"MARS FORCE MUSIC","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}]`
var collectionRightReportSingleOutput1 = `{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"15094dbb-9853-4737-aaa6-544ed27e0ac1","from":"PU200004","fromName":"MARS FORCE MUSIC","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}`
var collectionRightReportSingleOutput2 = `{"docType":"COPYRIGHTDATAREPORT","copyrightDataReportUUID":"2cfbdb47-cca7-3eca-b73e-0d6c478a6abc","isrc":"123Src","songTitle":"NY NY","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-11-15T22:27:34.111Z","rightHolders":[{"selector":"slct1","ipi":"ipi1","percent":42},{"selector":"slct2","ipi":"ipi2","percent":33}]}`
var updatedCollectionRightReportSingleInput = `[{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"15094dbb-9853-4737-aaa6-544ed27e0ac1","from":"PU200004","fromName":"MARS FORCE MUSIC - updated","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}]`
var updatedCollectionRightReportSingleOutput = `{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"15094dbb-9853-4737-aaa6-544ed27e0ac1","from":"PU200004","fromName":"MARS FORCE MUSIC - updated","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}`

var collectionRightReportMultipleOutput1 = `{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"15094dbb-9853-4737-aaa6-544ed27e0ac1","from":"PU200004","fromName":"MARS FORCE MUSIC","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}`
var collectionRightReportMultipleOutput2 = `{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"04240be9-73d3-4227-88a1-31c52d4db3bc","from":"PA300002","fromName":"URBAN SONGS","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}`

func MockGetCollectionRightReportResponse(functionName string) []byte {
	switch functionName {
	case "Test_AddCollectionRightReports_Single":
		return []byte(`{"successCount":1,"failureCount":0,"collectionRightsResponses":[]}`)
	case "Test_AddCollectionRightReports_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"collectionRightsResponses":[]}`)
	default:
		return []byte("[]")
	}
}

func MockGetCollectionRightReport(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	return []string{`{"docType":"COLLECTIONRIGHTREPORT","collectionRightUUID":"15094dbb-9853-4737-aaa6-544ed27e0ac1","from":"PU200004","fromName":"MARS FORCE MUSIC","startDate":"2010-12-1","endDate":"2030-12-1","rightHolders":[{"selector":"Territory=\"GER\"","ipi":"PG100001","percent":100},{"selector":"Territory=\"USA\"","ipi":"PU200001","percent":100},{"selector":"Territory=\"AUS\"","ipi":"PA300001","percent":100}]}`}, nil
}
func MockGetUpdatedCollectionRightReport(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	return []string{`{"docType":"COPYRIGHTDATAREPORT","copyrightDataReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a5abc","isrc":"1234567Src","songTitle":"modified","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-11-15T22:27:34.111Z","rightHolders":[{"selector":"slct1","ipi":"ipi1","percent":42},{"selector":"slct2","ipi":"ipi2","percent":33}]}`}, nil
}

func Test_AddCollectionRightReports_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getCollectionRightsForQueryString = MockGetCollectionRightReport
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addCollectionRights"), []byte(collectionRightReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	checkState(t, stub, collectionRightReportUUID, collectionRightReportSingleOutput1)

	expected := MockGetCollectionRightReportResponse("Test_AddCollectionRightReports_Single")
	// fmt.Println("-----------------------------")
	// fmt.Printf("actual - \n%s\n", actual)
	// fmt.Println("-----------------------------")
	// fmt.Printf("expected - \n%s\n", expected)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddCollectionRightReport_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addCollectionRights"), []byte(collectionRightReportMultipleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// // Check State for first  Report
	checkState(t, stub, collectionRightReportUUID, collectionRightReportMultipleOutput1)

	// // Check State for second  Report
	collectionRightReportUUID = "04240be9-73d3-4227-88a1-31c52d4db3bc"
	checkState(t, stub, collectionRightReportUUID, collectionRightReportMultipleOutput2)

	expected := MockGetCollectionRightReportResponse("Test_AddCollectionRightReports_Multiple")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_GetCollectionRightReportByID(t *testing.T) {

	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCollectionRights"), []byte(collectionRightReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	collectionRightReportUUID = "15094dbb-9853-4737-aaa6-544ed27e0ac1"
	actualReport, err := checkInvoke(t, stub, [][]byte{[]byte("getAssetByUUID"), []byte(collectionRightReportUUID)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedReports, err := MockGetCollectionRightReport(stub, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//normalize []string to string
	if !reflect.DeepEqual(expectedReports[0], string(actualReport)) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test_deleteCopyrightDataReportByIDs
func Test_DeleteCollectionRightReportByIDs(t *testing.T) {

	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getCollectionRightsForQueryString = MockGetCollectionRightReport
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCollectionRights"), []byte(collectionRightReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	//invoke delete method
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("deleteAssetByUUID"), []byte(collectionRightReportUUID)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	//we're expecting to delete a single record
	expected := `{"status":"200","message":"deleteAssetByUUID - deleted 1 records."}`
	if !reflect.DeepEqual(expected, string(actual)) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_updateCollectionRightReports_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Owner Administration
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCollectionRights"), []byte(collectionRightReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateCollectionRights"), []byte(updatedCollectionRightReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	collectionRightReportUUID = "15094dbb-9853-4737-aaa6-544ed27e0ac1"
	checkState(t, stub, collectionRightReportUUID, updatedCollectionRightReportSingleOutput)

	expected := MockGetCollectionRightReportResponse("Test_AddCollectionRightReports_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_AddCollectionRightReports_Empty1(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCollectionRights"), []byte(`[]`)})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_AddCollectionRightReports_Empty2(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCollectionRights"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}
}
