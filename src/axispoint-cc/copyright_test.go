package main

import (
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

const copyrightDataReportUUID = "1cfbdb47-cca7-3eca-b73e-0d6c478a5abc"

//var copyrightDataReportSingleInput = `[{"one":"1","two":"2","three":"3","docType":"COPYRIGHTDATAREPORT","copyrightDataReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eaa"}]`
var copyrightDataReportSingleInput = `[{"copyrightDataReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a5abc","docType":"COPYRIGHTDATAREPORT","isrc":"123Src","songTitle":"NY NY","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-11-15T22:27:34.111Z","rightHolders":[{"ipi":"ipi1","percent":42},{"ipi":"ipi2","percent":33}]}]`
var copyrightDataReportSingleOutput = `{"docType":"COPYRIGHTDATAREPORT","copyrightDataReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a5abc","isrc":"123Src","songTitle":"NY NY","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-11-15T22:27:34.111Z","rightHolders":[{"ipi":"ipi1","percent":42},{"ipi":"ipi2","percent":33}]}`
var updatedCopyrightDateReportSingleInput = `[{"copyrightDataReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a5abc","docType":"COPYRIGHTDATAREPORT","isrc":"1234567Src","songTitle":"modified","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-11-15T22:27:34.111Z","rightHolders":[{"ipi":"ipi1","percent":42},{"ipi":"ipi2","percent":33}]}]`

func MockGetCopyrightDataReportResponse(functionName string) []byte {
	switch functionName {
	case "Test_AddCopyrightDataReports_Single":
		return []byte(`{"successCount":1,"failureCount":0,"copyrightDataReports":[]}`)
	case "Test_AddCopyrightDataReports_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"copyrightDataReports":[]}`)
	default:
		return []byte("[]")
	}
}

func MockGetCopyrightDataReport(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	return []string{`{"docType":"COPYRIGHTDATAREPORT","copyrightDataReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a5abc","isrc":"123Src","songTitle":"NY NY","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-11-15T22:27:34.111Z","rightHolders":[{"ipi":"ipi1","percent":42},{"ipi":"ipi2","percent":33}]}`}, nil
}
func MockGetUpdatedCopyrightDataReport(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	return []string{`{"docType":"COPYRIGHTDATAREPORT","copyrightDataReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a5abc","isrc":"1234567Src","songTitle":"modified","startDate":"2018-01-01T21:17:34.371Z","endDate":"2018-11-15T22:27:34.111Z","rightHolders":[{"ipi":"ipi1","percent":42},{"ipi":"ipi2","percent":33}]}`}, nil
}
func Test_AddCopyrightDataReports_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getCopyrightDataReportForQueryString = MockGetCopyrightDataReport
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addCopyrightDataReports"), []byte(copyrightDataReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	checkState(t, stub, copyrightDataReportUUID, copyrightDataReportSingleOutput)

	expected := MockGetCopyrightDataReportResponse("Test_AddCopyrightDataReports_Single")
	// fmt.Println("-----------------------------")
	// fmt.Printf("actual - \n%s\n", actual)
	// fmt.Println("-----------------------------")
	// fmt.Printf("expected - \n%s\n", expected)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_GetCopyrightDataReportByID(t *testing.T) {

	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getCopyrightDataReportForQueryString = MockGetCopyrightDataReport
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCopyrightDataReports"), []byte(copyrightDataReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	getCopyrightDataReportForQueryString = MockGetCopyrightDataReport
	actualReport, err := checkInvoke(t, stub, [][]byte{[]byte("getCopyrightDataReportByIDs"), []byte(copyrightDataReportUUID)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedReports, err := MockGetCopyrightDataReport(stub, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//normalize []string to string
	if !reflect.DeepEqual("["+expectedReports[0]+"]", string(actualReport)) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test_deleteCopyrightDataReportByIDs
func Test_DeleteCopyrightDataReportByIDs(t *testing.T) {

	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getCopyrightDataReportForQueryString = MockGetCopyrightDataReport
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCopyrightDataReports"), []byte(copyrightDataReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	//invoke delete method
	actual, err := checkInvoke(t, stub, [][]byte{[]byte("deleteCopyrightDataReportByIDs"), []byte(copyrightDataReportUUID)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	//we're expecting to delete a single record
	expected := "deleted 1 records."

	if !reflect.DeepEqual(expected, string(actual)) {
		t.Fatalf("Actual response is not equal to expected response")
	}

}

//test - updateCopyrightDataReport
func Test_updateCopyrightDataReportByIDs(t *testing.T) {

	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	getCopyrightDataReportForQueryString = MockGetCopyrightDataReport
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addCopyrightDataReports"), []byte(copyrightDataReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	_, err = checkInvoke(t, stub, [][]byte{[]byte("updateCopyrightDataReports"), []byte(updatedCopyrightDateReportSingleInput)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	getCopyrightDataReportForQueryString = MockGetUpdatedCopyrightDataReport
	actualReport, err := checkInvoke(t, stub, [][]byte{[]byte("getCopyrightDataReportByIDs"), []byte(copyrightDataReportUUID)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedReports, err := MockGetUpdatedCopyrightDataReport(stub, "")
	if err != nil {
		t.Fatalf(err.Error())
	}
	//normalize []string to string
	if !reflect.DeepEqual("["+expectedReports[0]+"]", string(actualReport)) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}
