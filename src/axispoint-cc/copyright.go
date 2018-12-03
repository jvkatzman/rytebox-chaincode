package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var getCopyrightDataReportForQueryString = getObjectByQueryFromLedger

// addCopyrightDataReport - remove all data from the world state
// ================================================================================
func addCopyrightDataReport(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addCopyrightDataReport"
	logger.Info("ENTERING >", methodName, args)

	type CopyrightDataReportResponse struct {
		CopyrightDataReportUUID string `json:"copyrightDataReportUUID"`
		Message                 string `json:"message"`
		Success                 bool   `json:"success"`
	}

	type CopyrightDataReportOutput struct {
		SuccessCount   int                           `json:"successCount"`
		FailureCount   int                           `json:"failureCount"`
		RoyaltyReports []CopyrightDataReportResponse `json:"copyrightDataReports"`
	}

	if len(args) != 1 {
		return getErrorResponse("Missing arguments: Needed RoyaltyReport object to Create")
	}

	copyrightDataReportOutput := CopyrightDataReportOutput{}
	copyrightDataReports := &[]CopyrightDataReport{}
	copyrightDataReportResponses := []CopyrightDataReportResponse{}

	err := jsonToObject([]byte(args[0]), copyrightDataReports)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Royalty Reports
	for _, copyrightDataReport := range *copyrightDataReports {
		copyrightDataReport.DocType = COPYRIGHTDATAREPORT
		copyrightDataReportResponse := CopyrightDataReportResponse{}
		copyrightDataReportResponse.CopyrightDataReportUUID = copyrightDataReport.CopyrightDataUUID
		copyrightDataReportResponse.Success = true

		copyrightDataReportUUID, err := getCopyrightDataReportUUID(stub, copyrightDataReport)
		if err != nil {
			copyrightDataReportResponse.Success = false
			copyrightDataReportResponse.Message = err.Error()
			copyrightDataReportResponses = append(copyrightDataReportResponses, copyrightDataReportResponse)
			copyrightDataReportOutput.FailureCount++
			continue
		}

		copyrightDataReport.CopyrightDataUUID = copyrightDataReportUUID

		//Record royaltyReport on ledger
		copyrightDataReporBytes, err := objectToJSON(copyrightDataReport)
		if err != nil {
			copyrightDataReportResponse.Success = false
			copyrightDataReportResponse.Message = err.Error()
			copyrightDataReportResponses = append(copyrightDataReportResponses, copyrightDataReportResponse)
			copyrightDataReportOutput.FailureCount++
			continue
		}

		err = stub.PutState(copyrightDataReport.CopyrightDataUUID, copyrightDataReporBytes)
		if err != nil {
			copyrightDataReportResponse.Success = false
			copyrightDataReportResponse.Message = err.Error()
		}

		if copyrightDataReportResponse.Success {
			copyrightDataReportOutput.SuccessCount++
		} else {
			copyrightDataReportResponses = append(copyrightDataReportResponses, copyrightDataReportResponse)
			copyrightDataReportOutput.FailureCount++
		}
	}

	copyrightDataReportOutput.RoyaltyReports = copyrightDataReportResponses

	objBytes, _ := objectToJSON(copyrightDataReportOutput)
	logger.Info("EXITING <", methodName, copyrightDataReportOutput)
	return shim.Success(objBytes)
}

//getExploitationReportUUID : Get the UUID of the exploitation report based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
func getCopyrightDataReportUUID(stub shim.ChaincodeStubInterface, copyrightDataReport CopyrightDataReport) (string, error) {
	var methodName = "getCopyrightDataReportUUID"
	logger.Info("ENTERING >", methodName)

	getCopyrightDataReportUUID := ""
	//queryString := "{\"selector\":{\"docType\":\"" + EXPLOITATIONREPORT + "\",\"source\": \"" + royaltyReport.Source + "\",\"isrc\": \"" + royaltyReport.Isrc + "\",\"exploitationDate\": \"" + royaltyReport.ExploitationDate + "\",\"territory\": \"" + royaltyReport.Territory + "\",\"usageType\": \"" + royaltyReport.UsageType + "\"}}"
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"id\":\"%s\"}}", COPYRIGHTDATAREPORT, "????")
	logger.Info(methodName, queryString)

	queryResults, err := getCopyrightDataReportForQueryString(stub, queryString)
	if err != nil {
		return getCopyrightDataReportUUID, err
	}

	var copyrightDataReports []CopyrightDataReport
	err = sliceToStruct(queryResults, &copyrightDataReports)
	if err != nil {
		return getCopyrightDataReportUUID, err
	}

	if len(copyrightDataReports) <= 0 {
		//update message
		errorMessage := "UPDATE MESSAGE" //fmt.Sprintf("Cannot find Copyright Data Report with Source: %s, ISRC: %s, Exploitation Date: %s, Territory: %s, Usage Type: %s", royaltyReport.Source, royaltyReport.Isrc, royaltyReport.ExploitationDate, royaltyReport.Territory, royaltyReport.UsageType)
		return getCopyrightDataReportUUID, errors.New(errorMessage)
	}

	getCopyrightDataReportUUID = copyrightDataReports[0].CopyrightDataUUID

	logger.Info("EXITING <", methodName, getCopyrightDataReportUUID)
	return getCopyrightDataReportUUID, nil
}

// getCopyrightData - remove all data from the world state
// ================================================================================
func getCopyrightData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "getCopyrightData"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		message := fmt.Sprintf("Incorrect number of parameters received.")
		return shim.Error(message)
	}

	copyrightDataReportID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"id\":\"%s\"}}", COPYRIGHTDATAREPORT, copyrightDataReportID)
	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getCopyrightDataReportForQueryString(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	resultCopyrightReport := CopyrightDataReport{}
	err = sliceToStruct(queryResult, &resultCopyrightReport)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	copyrightReportResultBytes, err := objectToJSON(resultCopyrightReport)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Debugf("result(s) received from couch db: %s", string(copyrightReportResultBytes))

	return shim.Success(copyrightReportResultBytes)
}
