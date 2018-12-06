package main

import (
	"fmt"
	"testing"
)

// ********************************* Mock Data *********************************
var exploitationReport_in = `{"source":"M86321","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","isrc":"00029521","units":156062,"exploitationDate":"201811","amount":"36518.51","usageType":"SDIGM","exploitationReportUUID":"1cfbdb47-cca7-3eca-b73e-0d6c478a4eff","territory":"AUS"}`
var royaltyReport_in = `{"royaltyReportUUID":"a4c7408b-d68b-499e-8dfa-ff81b43ca8fe","source":"M86321","isrc":"00029521","songTitle":"HOLD THE LINE","writerName":"DAVID PAICH","units":156062,"exploitationDate":"20170131","amount":"7341.31000000","rightType":"SMECH","territory":"AUS","usageType":"SDIGM","target":"M86322"}`

// used for positive testing
var correctSimpleTrueExploitationSelector = "Source == 'M86321'"
var correctSimpleFalseExploitationSelector = "Isrc == '00000000'"
var correctTrueExploitationSelectorWithMultipleConjunctions = "Source == 'M86321' && 0 < Units && Units < 200000 && UsageType in ('FOO', 'SDIGM', 'BAR')"
var correctFalseSelectorWithMultipleConjunctions = "Source == 'M86321' && 0 < Units && Units < 200000 && UsageType in ('FOO', 'BAR', 'BAZ')"
var correctTrueSelectorWithMultipleLogicalOperators = "(WriterName == 'DAVID PAICH' || Isrc == '00029521') && (SongTitle == 'HOLD THE LINE' || Units > 200000)"
var correctTrueRoyaltySelector = "Source == 'M86321' && RightType == 'SMECH' && Territory in ('AUS', 'USA', 'GBR')"

// used for negative testing
var simpleMissingFieldSelectorWith = "Foo == 'bar'"
var randomString = "fooBarBaz"

// *****************************************************************************

// ****************************** Positive tests *******************************
func TestEvaluate_CorrectSimpleSelectorsForExploitation_ShouldReturnTrueThenFalse(t *testing.T) {
	var exploitationReport = getExploitationReport(exploitationReport_in, t)
	testEval(correctSimpleTrueExploitationSelector, &exploitationReport, false, true, t)
	testEval(correctSimpleFalseExploitationSelector, &exploitationReport, false, false, t)
}

func TestEvaluate_CorrectSelectorMultipleConjunctionForExploitation_ShouldReturnTrueThenFalse(t *testing.T) {
	var exploitationReport = getExploitationReport(exploitationReport_in, t)
	testEval(correctTrueExploitationSelectorWithMultipleConjunctions, &exploitationReport, false, true, t)
	testEval(correctFalseSelectorWithMultipleConjunctions, &exploitationReport, false, true, t)
}

func TestEvaluate_CorrectSelectorMultipleLogicalOperators_ShouldReturnTrue(t *testing.T) {
	var exploitationReport = getExploitationReport(exploitationReport_in, t)
	testEval(correctTrueSelectorWithMultipleLogicalOperators, &exploitationReport, false, true, t)
}

func TestEvaluate_CorrectTrueRoyaltySelector_ShouldReturnTrue(t *testing.T) {
	var royaltyReport = getRoyaltyReport(royaltyReport_in, t)
	testEval(correctTrueRoyaltySelector, &royaltyReport, false, true, t)
}

// ****************************** Negative tests ******************************
func TestEvaluate_CorrectSelectorContainingExploitationMissingFiled_ShouldReturnError(t *testing.T) {
	var exploitationReport = getExploitationReport(exploitationReport_in, t)
	testEval(simpleMissingFieldSelectorWith, &exploitationReport, true, nil, t)
}

func TestEvaluate_WrongInputs_ShouldReturnError(t *testing.T) {
	var royaltyReport = getRoyaltyReport(royaltyReport_in, t)

	// test for wrong input asset type (struct instead of a struct pointer)
	testEval(simpleMissingFieldSelectorWith, royaltyReport, true, nil, t)

	//// test for wrong input asset type (string pointer instead of a struct pointer)
	testEval(simpleMissingFieldSelectorWith, &randomString, true, nil, t)

	// test for wrong input asset type (nil instead of a struct)
	testEval(simpleMissingFieldSelectorWith, nil, true, nil, t)

	// test for wrong input selector (random string)
	testEval(randomString, &royaltyReport, true, nil, t)

	// test for wrong input selector (empty string)
	testEval("", &royaltyReport, true, nil, t)

	// test for wrong input selector (empty string)
	testEval("", &royaltyReport, true, nil, t)
}

// ****************************** Utils for tests ******************************
func getExploitationReport(exploitationReportIn string, t *testing.T) ExploitationReport {
	var exploitationReport = ExploitationReport{}
	err := jsonToObject([]byte(exploitationReportIn), &exploitationReport)

	if err != nil {
		t.Fatalf("Failed to parse input JSON: %s", err.Error())
	}

	return exploitationReport
}

func getRoyaltyReport(royaltyReportIn string, t *testing.T) RoyaltyReport {
	var royaltyReport = RoyaltyReport{}
	err := jsonToObject([]byte(royaltyReportIn), &royaltyReport)

	if err != nil {
		t.Fatalf("Failed to parse input JSON: %s", err.Error())
	}

	return royaltyReport
}

func testEval(selector string, asset interface{}, isErrExpected bool, expectedResult interface{}, t *testing.T) {
	if !isErrExpected && expectedResult == nil {
		t.Fatalf("The expected result cannot be nil if no err is expected")
	}

	result, err := evaluate(selector, asset)
	fmt.Println("Evaluate result =", result)
	fmt.Println("Evaluate error =", err)

	if isErrExpected {
		if err == nil || result != nil {
			t.Fatalf("Eval should have returned an error")
		}
	} else {
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if result != expectedResult {
			t.Fatalf("Evaluate returned %t, expected %t", result, expectedResult)
		}
	}
}
