package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var getCopyrightDataReportForQueryString = getObjectByQueryFromLedger

// addCopyrightDataReport - add a copyright data report
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
		//UUID is already present is this valid???
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

// getCopyrightDataByID - retrieve a copyright data report by id by an array
// ================================================================================
func getCopyrightDataReportByIDs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "getCopyrightDataReportByIDs"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		message := fmt.Sprintf("%s - Incorrect number of parameters received.", methodName)
		logger.Error(message)
		return shim.Error(message)
	}
	inSubQuery := `{"$in":[`

	for _, copyrightDataReportUUID := range args {
		inSubQuery += fmt.Sprintf("\"%s\",", copyrightDataReportUUID)
	}

	//remove the last commma and add the remaining closing tags
	inSubQuery = strings.TrimSuffix(inSubQuery, ",") + "]}"
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"copyrightDataReportUUID\":%s}}", COPYRIGHTDATAREPORT, inSubQuery)
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
	copyrightReportResultBytes, err := objectToJSON(resultCopyrightReports)
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
func updateCopyrightDataReports(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "updateCopyrightDataReports"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		message := fmt.Sprintf("%s - Incorrect number of parameters received.", methodName)
		logger.Error(message)
		return shim.Error(message)
	}
	updatedCopyrightDataReports := &[]CopyrightDataReport{}

	err := jsonToObject([]byte(args[0]), updatedCopyrightDataReports)
	if err != nil {
		logger.Errorf("%s - failed to convert ")
		return getErrorResponse(err.Error())
	}

	for _, updatedCopyrightDataReport := range *updatedCopyrightDataReports {
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
	}

	return shim.Success(nil)
}

// searchForCopyrightDataReportWithParameters - search for copyright data report(s)
// method expects an argument list where
// args[0] must be 'isrc'
// args[1] must be 'songTitle'
// args[2] must be 'startDate'
// args[3] must be 'endDate'
// ================================================================================
func searchForCopyrightDataReportWithParameters(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "searchForCopyrightDataReportWithParameters"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		message := fmt.Sprintf("%s - incorrect # of arguments received.", methodName)
		logger.Error(message)
		return shim.Error(message)
	}
	var queryString string
	//expected arguments
	switch len(args) {
	case 1: //isrc
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"isrc\":\"%s\"}}", COPYRIGHTDATAREPORT, args[0])
	case 2: //isrc && song title
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"isrc\":\"%s\", \"songTitle\":\"%s\"}}", COPYRIGHTDATAREPORT, args[0], args[1])
	case 3: //isrc && song title && start date
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"isrc\":\"%s\", \"songTitle\":\"%s\", \"startDate\":\"%s\"}}", COPYRIGHTDATAREPORT, args[0], args[1], args[2])
	case 4: ////isrc && song title && start and end dates
		queryString = fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"isrc\":\"%s\", \"songTitle\":\"%s\", \"startDate\":\"%s\", \"endDate\":\"%s\"}}", COPYRIGHTDATAREPORT, args[0], args[1], args[2], args[3])
	default:
		errMsg := fmt.Sprintf("%s - Failed to determine provided args length. arguments : '%s'.", methodName, strings.Join(args, ","))
		logger.Errorf(errMsg)
		return shim.Error(errMsg)
	}
	logger.Infof("%s - executing couch db query : %s", methodName, queryString)
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
	copyrightReportResultBytes, err := objectToJSON(resultCopyrightReports)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Debugf("result(s) received from couch db: %s", string(copyrightReportResultBytes))

	return shim.Success(copyrightReportResultBytes)
}

//getAllCopyrightDataReports: get all the copyright data reports that exist
func getAllCopyrightDataReports(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getAllCopyrightDataReports"
	logger.Infof("%s - Begin Execution ", methodName)
	defer logger.Infof("%s - End Execution ", methodName)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", COPYRIGHTDATAREPORT)
	if len(args) > 1 {
		queryString = args[0]
	}

	logger.Info("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getCopyrightDataReportForQueryString(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	var resultCopyrightDataReports []CopyrightDataReport
	err = sliceToStruct(queryResult, &resultCopyrightDataReports)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	queryResultBytes, err := objectToJSON(resultCopyrightDataReports)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Info("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}
