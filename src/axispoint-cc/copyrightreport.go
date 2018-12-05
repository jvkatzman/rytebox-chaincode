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
func addCopyrightDataReports(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addCopyrightDataReports"
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

	// Iterate over Copyright Reports
	for _, copyrightDataReport := range *copyrightDataReports {
		copyrightDataReport.DocType = COPYRIGHTDATAREPORT
		copyrightDataReportResponse := CopyrightDataReportResponse{}
		copyrightDataReportResponse.CopyrightDataReportUUID = copyrightDataReport.CopyrightDataUUID
		copyrightDataReportResponse.Success = true

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

	//return "1cfbdb47-cca7-3eca-b73e-0d6c478a4eaa", nil

	getCopyrightDataReportUUID := ""
	//queryString := "{\"selector\":{\"docType\":\"" + EXPLOITATIONREPORT + "\",\"source\": \"" + royaltyReport.Source + "\",\"isrc\": \"" + royaltyReport.Isrc + "\",\"exploitationDate\": \"" + royaltyReport.ExploitationDate + "\",\"territory\": \"" + royaltyReport.Territory + "\",\"usageType\": \"" + royaltyReport.UsageType + "\"}}"
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"isrc\":\"%s\"}}", COPYRIGHTDATAREPORT, copyrightDataReport.Isrc)
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
		errorMessage := fmt.Sprintf("Cannot find Copyright Data Report with ISRC: %s", copyrightDataReport.Isrc)
		return getCopyrightDataReportUUID, errors.New(errorMessage)
	}

	getCopyrightDataReportUUID = copyrightDataReports[0].CopyrightDataUUID

	logger.Info("EXITING <", methodName, getCopyrightDataReportUUID)
	return getCopyrightDataReportUUID, nil
}

// getCopyrightDataByID - retrieve a copyright data report by id by an array
// ================================================================================
func getCopyrightDataReportByID(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "getCopyrightDataReportByID"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		message := fmt.Sprintf("%s - Incorrect number of parameters received.", methodName)
		return shim.Error(message)
	}

	copyrightDataReportUUID := args[0]
	copyrightDataReportIsrc := ""

	//queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"copyrightDataReportUUID\":\"%s\"}}", COPYRIGHTDATAREPORT, copyrightDataReportID)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"$or\":[{\"copyrightDataReportUUID\":\"%s\"},{\"isrc\":\"%s\"}]}}", COPYRIGHTDATAREPORT, copyrightDataReportUUID, copyrightDataReportIsrc)
	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getCopyrightDataReportForQueryString(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	resultCopyrightReports := []CopyrightDataReport{}
	err = sliceToStruct(queryResult, &resultCopyrightReports)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// we should just have a single item in the result array
	copyrightReportResultBytes, err := objectToJSON(resultCopyrightReports[0])
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Debugf("result(s) received from couch db: %s", string(copyrightReportResultBytes))

	return shim.Success(copyrightReportResultBytes)
}

// deleteCopyrightDataByIDs - delete a copyright data report by ids in an array
// ================================================================================
func deleteCopyrightDataReportByIDs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "deleteCopyrightDataReportByIDs"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)
	deletedRecordCount := 0

	if len(args) < 1 {
		message := fmt.Sprintf("%s - Incorrect number of parameters received.", methodName)
		logger.Error(message)
		return shim.Error(message)
	}

	for _, copyrightDataReportUUID := range args {
		logger.Infof("%s - deleting copyright record with uuid: %s", methodName, copyrightDataReportUUID)
		err := stub.DelState(copyrightDataReportUUID)
		if err != nil {
			message := fmt.Sprintf("%s - Failed to delete copyright data report with id : %s", methodName, copyrightDataReportUUID)
			logger.Info(message)
		} else {
			deletedRecordCount++
		}
	}

	logger.Infof("%s - successfully deleted %d records.", methodName, deletedRecordCount)
	return shim.Success([]byte(fmt.Sprintf("deleted %d records.", deletedRecordCount)))
}

// updateCopyrightDataReport - update an existing copyright data report
// ================================================================================
func updateCopyrightDataReport(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "updateCopyrightDataReport"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		message := fmt.Sprintf("%s - Incorrect number of parameters received.", methodName)
		logger.Error(message)
		return shim.Error(message)
	}
	updatedCopyrightDataReport := CopyrightDataReport{}

	err := jsonToObject([]byte(args[0]), &updatedCopyrightDataReport)
	if err != nil {
		logger.Errorf("%s - failed to convert ")
		return getErrorResponse(err.Error())
	}

	existingReportBytes, err := stub.GetState(updatedCopyrightDataReport.CopyrightDataUUID)
	if err != nil {
		message := fmt.Sprintf("%s - Failed to check if the existing report with id %s can be updated.", methodName, updatedCopyrightDataReport.CopyrightDataUUID)
		logger.Error(message)
		return shim.Error(message)
	}

	if existingReportBytes == nil {
		message := fmt.Sprintf("%s - report with id %s cannot be updated since it was not found on the ledger.", methodName, updatedCopyrightDataReport.CopyrightDataUUID)
		logger.Error(message)
		return shim.Error(message)
	}

	updatedReportBytes, err := objectToJSON(updatedCopyrightDataReport)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	err = stub.PutState(updatedCopyrightDataReport.CopyrightDataUUID, updatedReportBytes)
	if err != nil {
		message := fmt.Sprintf("%s - Failed to update existing report with id %s wth error : %s.", methodName, updatedCopyrightDataReport.CopyrightDataUUID, err.Error())
		logger.Error(message)
		return shim.Error(message)
	}

	logger.Infof("%s - successfully updated existing report with id: %s", methodName, updatedCopyrightDataReport.CopyrightDataUUID)

	return shim.Success(nil)
}
