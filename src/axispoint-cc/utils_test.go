package main

import (
	"fmt"
	"testing"
)

// ********************************* Mock Data *********************************
var exploitationReport_in = `{ "source": "M86321", "songTitle": "LIVING WITH THE LAW", "writerName": "CHRIS WHITLEY", "isrc": "00055521", "units": 456, "exploitationDate": "20170131", "amount": "69.71000000", "usageType": "SDIGM", "exploitationReportUUID": "b6d7a629-85c5-36a6-96fa-8dc3a5f71169", "territory": "AUS" }`
var royaltyReport_in = `{ "territory": "AUS", "songTitle": "GECKOS!!", "writerName": "KIERAN CASH", "isrc": "00055524", "units": 140, "exploitationDate": "20170131", "amount": "14.094", "usageType": "SDIGP", "source": "M86321", "representative": "PP8819H", "collector": "PP8819H", "rightHolder": "W998", "royaltyStatementUUID": "31f52320-d090-3bc1-935a-b2bc9becb6cb", "rightType": "PERF" }`

// used for positive testing
var correctSimpleTrueExploitationSelector = "Source == 'M86321'"
var correctSimpleFalseExploitationSelector = "Isrc == '00000000'"
var correctTrueExploitationSelectorWithMultipleConjunctions = "Source == 'M86321' && 0 < Units && Units < 200000 && UsageType in ('FOO', 'SDIGM', 'BAR')"
var correctFalseSelectorWithMultipleConjunctions = "Source == 'M86321' && 0 < Units && Units < 200000 && UsageType in ('FOO', 'BAR', 'BAZ')"
var correctTrueSelectorWithMultipleLogicalOperators = "(WriterName == 'CHRIS WHITLEY' || Isrc == '00055521') && (SongTitle == 'LIVING WITH THE LAW' || Units > 200000)"
var correctTrueRoyaltySelector = "Source == 'M86321' && RightType == 'PERF' && Territory in ('AUS', 'USA', 'GBR')"

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
	testEval(correctFalseSelectorWithMultipleConjunctions, &exploitationReport, false, false, t)
}

func TestEvaluate_CorrectSelectorMultipleLogicalOperators_ShouldReturnTrue(t *testing.T) {
	var exploitationReport = getExploitationReport(exploitationReport_in, t)
	testEval(correctTrueSelectorWithMultipleLogicalOperators, &exploitationReport, false, true, t)
}

func TestEvaluate_CorrectTrueRoyaltySelector_ShouldReturnTrue(t *testing.T) {
	var royaltyReport = getRoyaltyStatement(royaltyReport_in, t)
	testEval(correctTrueRoyaltySelector, &royaltyReport, false, true, t)
}

// ****************************** Negative tests ******************************
func TestEvaluate_CorrectSelectorContainingExploitationMissingFiled_ShouldReturnError(t *testing.T) {
	var exploitationReport = getExploitationReport(exploitationReport_in, t)
	testEval(simpleMissingFieldSelectorWith, &exploitationReport, true, nil, t)
}

func TestEvaluate_WrongInputs_ShouldReturnError(t *testing.T) {
	var royaltyReport = getRoyaltyStatement(royaltyReport_in, t)

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

func getRoyaltyStatement(royaltyReportIn string, t *testing.T) RoyaltyStatement {
	var royaltyReport = RoyaltyStatement{}
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
