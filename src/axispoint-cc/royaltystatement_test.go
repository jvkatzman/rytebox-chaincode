package main

import (
	"reflect"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var royaltyStatementSingle1_in = `[{"royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":200,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Ned-IPI","administrator":"Swedish-Publishing-IPI","collector":"","state":"MISSING_AFFILIATE"}]`
var royaltyStatementSingle1_out = `{"docType":"ROYALTYSTATEMENT","royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":200,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Ned-IPI","administrator":"Swedish-Publishing-IPI","collector":"","state":"MISSING_AFFILIATE"}`
var royaltyStatementMultiple1_in = `[{"royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":200,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Ned-IPI","administrator":"Swedish-Publishing-IPI","collector":"","state":"MISSING_AFFILIATE"},{"royaltyStatementUUID":"5bbbda3a-6335-4248-9d10-019a73f59dfc","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":300,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Homer-Simpson-IPI","administrator":"ACME-Music-Corp-IPI","collector":"","state":"MISSING_AFFILIATE"}]`
var royaltyStatementSingle2_out = `{"docType":"ROYALTYSTATEMENT","royaltyStatementUUID":"5bbbda3a-6335-4248-9d10-019a73f59dfc","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":300,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Homer-Simpson-IPI","administrator":"ACME-Music-Corp-IPI","collector":"","state":"MISSING_AFFILIATE"}`
var royaltyStatementSingle2_in = `[{"royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":200,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Ned-IPI","administrator":"Swedish-Publishing-IPI","collector":"","state":"MISSING_AFFILIATE"}]`
var royaltyStatementMultiple2_in = `[{"royaltyStatementUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","source":"M86321","isrc":"00029524","exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"},{"royaltyStatementUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8ff","source":"M86321","isrc":"00029525","exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}]`

func MockGetRoyaltyStatementResponse(functionName string) []byte {
	switch functionName {
	case "Test_addRoyaltyStatements_Single":
		return []byte(`{"successCount":1,"failureCount":0,"royaltyStatements":[]}`)
	case "Test_addRoyaltyStatements_Single_Failure":
		return []byte(`{"successCount":0,"failureCount":1,"royaltyStatements":[{"royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","message":"Royalty Statement already exists!","success":false}]}`)
	case "Test_addRoyaltyStatements_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"royaltyStatements":[]}`)
	case "Test_addRoyaltyStatements_Multiple_Failure":
		return []byte(`{"successCount":0,"failureCount":2,"royaltyStatements":[{"royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","message":"Royalty Statement already exists!","success":false},{"royaltyStatementUUID":"5bbbda3a-6335-4248-9d10-019a73f59dfc","message":"Royalty Statement already exists!","success":false}]}`)
	case "Test_GetRoyaltyStatements":
		return []byte(`[{"docType":"ROYALTYSTATEMENT","royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":200,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Ned-IPI","administrator":"Swedish-Publishing-IPI","collector":"","state":"MISSING_AFFILIATE"}]`)
	case "Test_GetRoyaltyStatementByUUID":
		return []byte(`{"docType":"ROYALTYSTATEMENT","royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":200,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Ned-IPI","administrator":"Swedish-Publishing-IPI","collector":"","state":"MISSING_AFFILIATE"}`)
	case "Test_GetRoyaltyStatementByUUID_Failure":
		return []byte(`{"status":"500","message":"UUID: 85fff2bf-00a2-423b-9567-55c6f4ee6ee2 does not exist"}`)
	case "Test_UpdateRoyaltyStatements_Single":
		return []byte(`{"successCount":1,"failureCount":0,"royaltyStatements":[]}`)
	case "Test_UpdateRoyaltyStatements_Multiple":
		return []byte(`{"successCount":2,"failureCount":0,"royaltyStatements":[]}`)
	default:
		return []byte("[]")
	}
}

func MockGetRoyaltyStatementQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]string, error) {
	//return []string{`{"amount":"18.5034375","docType":"ROYALTYSTATEMENT","exploitationDate":"20170131","exploitationReportUUID":"7cb134b2-156f-32e0-a4d9-6165c6ad1aca","isrc":"00055524","rightType":"PERF","royaltyStatementUUID":"2e0d3d18-f25f-36d0-81ef-6e3da893d4aa","songTitle":"GECKOS!!","source":"P8819H","target":"W998","territory":"AUS","units":226,"usageType":"SDIGP","writerName":"KIERAN CASH"},{"amount":"12.404250000000001","docType":"ROYALTYSTATEMENT","exploitationDate":"20170131","exploitationReportUUID":"6874280d-2897-3321-b238-0b4dfa0aa516","isrc":"00055524","rightType":"MECH","royaltyStatementUUID":"7f384cbf-0d0d-3698-9714-841b8ecb73f9","songTitle":"GECKOS!!","source":"P8819H","target":"W998","territory":"AUS","units":150,"usageType":"SDIGM","writerName":"KIERAN CASH"},{"amount":"12.366","docType":"ROYALTYSTATEMENT","exploitationDate":"20170131","exploitationReportUUID":"8ab33826-399f-3707-a0af-dfedc3d3b7f3","isrc":"00055524","rightType":"MECH","royaltyStatementUUID":"94c878c5-f754-3b04-b90e-4f01cbd54ad6","songTitle":"GECKOS!!","source":"P8819H","target":"W998","territory":"AUS","units":164,"usageType":"SMECH","writerName":"KIERAN CASH"}`}, nil
	return []string{`{"docType":"ROYALTYSTATEMENT","royaltyStatementUUID":"0daccfc9-9e3a-43f1-8e60-d0d0916a82e3","exploitationReportUUID":"5c42d472-d137-4218-94a3-2825999a6f10","source":"spotify-IPI","isrc":"FlowThroughSong1-ISRC","songTitle":"missing affiliation song","writerName":"Homer, Ned","units":10000,"exploitationDate":"2018-12-30T00:00:00.000Z","amount":200,"rightType":"OWNERSHIP","territory":"FRA","usageType":"MECH","rightHolder":"Ned-IPI","administrator":"Swedish-Publishing-IPI","collector":"","state":"MISSING_AFFILIATE"}`}, nil
}

func Test_addRoyaltyStatements_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var royalReportUUID = "0daccfc9-9e3a-43f1-8e60-d0d0916a82e3"
	checkState(t, stub, royalReportUUID, royaltyStatementSingle1_out)

	expected := MockGetRoyaltyStatementResponse("Test_addRoyaltyStatements_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_addRoyaltyStatements_Single_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementSingle2_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementSingle2_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetRoyaltyStatementResponse("Test_addRoyaltyStatements_Single_Failure")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_addRoyaltyStatements_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementMultiple1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var royalReportUUID = "0daccfc9-9e3a-43f1-8e60-d0d0916a82e3"
	checkState(t, stub, royalReportUUID, royaltyStatementSingle1_out)

	royalReportUUID = "5bbbda3a-6335-4248-9d10-019a73f59dfc"
	checkState(t, stub, royalReportUUID, royaltyStatementSingle2_out)

	expected := MockGetRoyaltyStatementResponse("Test_addRoyaltyStatements_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_addRoyaltyStatements_Multiple_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementMultiple1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementMultiple1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetRoyaltyStatementResponse("Test_addRoyaltyStatements_Multiple_Failure")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_AddRoyaltyStatements_Empty1(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(`[]`)})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_AddRoyaltyStatements_Empty2(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}
}

//Test the Edge cases
func Test_GetRoyaltyStatements(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	getRoyaltyStatementsForQueryString = MockGetRoyaltyStatementQueryResultForQueryString

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getRoyaltyStatements"), []byte("")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetRoyaltyStatementResponse("Test_GetRoyaltyStatements")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetRoyaltyStatementByUUID(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Administrator Affiliation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getAssetByUUID"), []byte("0daccfc9-9e3a-43f1-8e60-d0d0916a82e3")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetRoyaltyStatementResponse("Test_GetRoyaltyStatementByUUID")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

//Test the Edge cases
func Test_GetRoyaltyStatementByUUID_Failure(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Administrator Affiliation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("getAssetByUUID"), []byte("85fff2bf-00a2-423b-9567-55c6f4ee6ee2")})
	if err != nil {
		t.Fatalf(err.Error())
	}

	expected := MockGetRoyaltyStatementResponse("Test_GetRoyaltyStatementByUUID_Failure")

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateRoyaltyStatements_Single(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)
	// Add Administrator Affiliation
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateRoyaltyStatements"), []byte(royaltyStatementSingle1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for Transaction
	var royaltyStatementUUID = "0daccfc9-9e3a-43f1-8e60-d0d0916a82e3"
	checkState(t, stub, royaltyStatementUUID, royaltyStatementSingle1_out)

	expected := MockGetRoyaltyStatementResponse("Test_UpdateRoyaltyStatements_Single")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}

func Test_UpdateRoyaltyStatements_Multiple(t *testing.T) {
	scc := new(AxispointChaincode)
	stub := shim.NewMockStub("AxispointChaincode", scc)

	// Init
	checkInit(t, stub, [][]byte{[]byte("init"), []byte("")}, nil)

	// Add Administrator Affiliations
	_, err := checkInvoke(t, stub, [][]byte{[]byte("addRoyaltyStatements"), []byte(royaltyStatementMultiple1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	actual, err := checkInvoke(t, stub, [][]byte{[]byte("updateRoyaltyStatements"), []byte(royaltyStatementMultiple1_in)})
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Check State for first Exploitation Report
	var royalReportUUID = "0daccfc9-9e3a-43f1-8e60-d0d0916a82e3"
	checkState(t, stub, royalReportUUID, royaltyStatementSingle1_out)

	royalReportUUID = "5bbbda3a-6335-4248-9d10-019a73f59dfc"
	checkState(t, stub, royalReportUUID, royaltyStatementSingle2_out)

	expected := MockGetRoyaltyStatementResponse("Test_UpdateRoyaltyStatements_Multiple")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Actual response is not equal to expected response")
	}
}
