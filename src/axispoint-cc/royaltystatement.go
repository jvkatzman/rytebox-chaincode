/*
Copyright IBM Corp.. 2018 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//RoyaltyStatementResponse : defines response data from blockchain request
type RoyaltyStatementResponse struct {
	RoyaltyStatementUUID string `json:"royaltyStatementUUID"`
	Message              string `json:"message"`
	Success              bool   `json:"success"`
}

//RoyaltyStatementOutput : defines accumulated output of blockchain requests
type RoyaltyStatementOutput struct {
	SuccessCount     int                        `json:"successCount"`
	FailureCount     int                        `json:"failureCount"`
	RoyaltyStatement []RoyaltyStatementResponse `json:"royaltyStatement"`
}

// GetExploitationReportForQueryString : Get exploitation reports based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
var getExploitationReportForQueryString = getObjectByQueryFromLedger
var getRoyaltyStatementsForQueryString = getObjectByQueryFromLedger

//AddRoyaltyStatements : Add Royalty Statements to the ledger
func addRoyaltyStatements(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addRoyaltyStatements"
	logger.Info("ENTERING >", methodName, args)

	if len(args) != 1 {
		return getErrorResponse("Missing arguments: Needed RoyaltyStatement object to Create")
	}

	royaltyStatementOutput := RoyaltyStatementOutput{}
	royaltyStatements := &[]RoyaltyStatement{}
	royaltyStatementResponses := []RoyaltyStatementResponse{}

	err := jsonToObject([]byte(args[0]), royaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Royalty Statements
	for _, royaltyStatement := range *royaltyStatements {
		royaltyStatement.DocType = ROYALTYSTATEMENT
		royaltyStatementResponse := RoyaltyStatementResponse{}
		royaltyStatementResponse.RoyaltyStatementUUID = royaltyStatement.RoyaltyStatementUUID
		royaltyStatementResponse.Success = true

		exploitationReportUUID, err := getExploitationReportUUID(stub, royaltyStatement)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
			continue
		}

		royaltyStatement.ExploitationReportUUID = exploitationReportUUID

		//Record royaltyReport on ledger
		royaltyStatementBytes, err := objectToJSON(royaltyStatement)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
			continue
		}

		err = stub.PutState(royaltyStatement.RoyaltyStatementUUID, royaltyStatementBytes)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
		}

		if royaltyStatementResponse.Success {
			royaltyStatementOutput.SuccessCount++
		} else {
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
		}
	}

	royaltyStatementOutput.RoyaltyStatement = royaltyStatementResponses

	objBytes, _ := objectToJSON(royaltyStatementOutput)
	logger.Info("EXITING <", methodName, royaltyStatementOutput)
	return shim.Success(objBytes)
}

//getExploitationReportUUID : Get the UUID of the exploitation report based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
func getExploitationReportUUID(stub shim.ChaincodeStubInterface, royaltyReport RoyaltyStatement) (string, error) {
	var methodName = "getExploitationReportUUID"
	logger.Info("ENTERING >", methodName)

	exploitationReportUUID := ""
	queryString := "{\"selector\":{\"docType\":\"" + EXPLOITATIONREPORT + "\",\"source\": \"" + royaltyReport.Source + "\",\"isrc\": \"" + royaltyReport.Isrc + "\",\"exploitationDate\": \"" + royaltyReport.ExploitationDate + "\",\"territory\": \"" + royaltyReport.Territory + "\",\"usageType\": \"" + royaltyReport.UsageType + "\"}}"
	logger.Info(methodName, queryString)

	queryResults, err := getExploitationReportForQueryString(stub, queryString)
	if err != nil {
		return exploitationReportUUID, err
	}

	var exploitationReports []ExploitationReport
	err = sliceToStruct(queryResults, &exploitationReports)
	if err != nil {
		return exploitationReportUUID, err
	}

	if len(exploitationReports) <= 0 {
		errorMessage := fmt.Sprintf("Cannot find Exploitation Report with Source: %s, ISRC: %s, Exploitation Date: %s, Territory: %s, Usage Type: %s", royaltyReport.Source, royaltyReport.Isrc, royaltyReport.ExploitationDate, royaltyReport.Territory, royaltyReport.UsageType)
		return exploitationReportUUID, errors.New(errorMessage)
	}

	exploitationReportUUID = exploitationReports[0].ExploitationReportUUID

	logger.Info("EXITING <", methodName, exploitationReportUUID)
	return exploitationReportUUID, nil
}

//get paid royalty data for a given period
//expected parameters: exploitation date and the target(the creator)
func getRoyaltyDataForPeriod(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getRoyaltyDataForPeriod"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 2 {
		errMsg := fmt.Sprintf("%s - Incorrect number of parameters provided : %s.  Expecting exploitation date and target.", methodName, strings.Join(args, ","))
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}

	exploitationDate := args[0]
	targetCreator := args[1]
	//do a rich query to get the data from the ledger
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"target\":\"%s\",\"exploitationDate\":{\"$lte\":\"%s\"}}}", ROYALTYSTATEMENT, targetCreator, exploitationDate)
	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getRoyaltyStatementsForQueryString(stub, queryString) //getQueryResultInBytes(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	var resultRoyaltyStatements []RoyaltyStatement
	err = sliceToStruct(queryResult, &resultRoyaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	queryResultBytes, err := objectToJSON(resultRoyaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Infof("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}

//get paid royalty data based on a selector string
//expected parameters: selector string
func queryRoyaltyStatements(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "queryRoyaltyStatements"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		errMsg := fmt.Sprintf("%s - Incorrect number of parameters provided : %s.  Expecting selector string.", methodName, strings.Join(args, ","))
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}

	//do a rich query to get the data from the ledger
	queryString := args[0]
	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getRoyaltyStatementsForQueryString(stub, queryString) //getQueryResultInBytes(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	var resultRoyaltyStatements []RoyaltyStatement
	err = sliceToStruct(queryResult, &resultRoyaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	queryResultBytes, err := objectToJSON(resultRoyaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Infof("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}

// getRoyaltyStatementsByUUID - retrieve royalty statements by id in an array
// ================================================================================
func getRoyaltyStatementsByUUIDs(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "getRoyaltyStatementsByUUIDs"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 1 {
		message := fmt.Sprintf("%s - Incorrect number of parameters received.", methodName)
		logger.Error(message)
		return shim.Error(message)
	}
	inSubQuery := `{"$in":[`

	for _, royaltyStatementUUIDs := range args {
		inSubQuery += fmt.Sprintf("\"%s\",", royaltyStatementUUIDs)
	}

	//remove the last commma and add the remaining closing tags
	inSubQuery = strings.TrimSuffix(inSubQuery, ",") + "]}"
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"royaltyStatementUUID\":%s}}", ROYALTYSTATEMENT, inSubQuery)
	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getObjectByQueryFromLedger(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	resultRoyaltyStatements := []RoyaltyStatement{}
	err = sliceToStruct(queryResult, &resultRoyaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// we should just have a single item in the result array
	royaltyStatementResultBytes, err := objectToJSON(resultRoyaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	logger.Debugf("result(s) received from couch db: %s", string(royaltyStatementResultBytes))

	return shim.Success(royaltyStatementResultBytes)
}

// updateRoyaltyStatements function contains business logic to update
// Royalty Statements on the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of royalty statements.
* @return   {pb.Response}    - peer Response
 */
func updateRoyaltyStatements(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "updateRoyaltyStatements"
	logger.Info("ENTERING >", methodName, args)

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Royalty Statement objects is required")
	}

	royaltyStatementsOutput := RoyaltyStatementOutput{}
	royaltyStatements := &[]RoyaltyStatement{}
	royaltyStatementResponses := []RoyaltyStatementResponse{}

	//Unmarshal the args input to an array of owner administration records
	err := jsonToObject([]byte(args[0]), royaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Exploitation Reports
	for _, royaltyStatement := range *royaltyStatements {
		royaltyStatement.DocType = ROYALTYSTATEMENT
		royaltyStatementResponse := RoyaltyStatementResponse{}
		royaltyStatementResponse.RoyaltyStatementUUID = royaltyStatement.RoyaltyStatementUUID
		royaltyStatementResponse.Success = true

		//Record Exploitation Report on ledger
		royaltyStatementsBytes, err := objectToJSON(royaltyStatement)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementsOutput.FailureCount++
			continue
		}

		royaltyStatementExistingBytes, err := stub.GetState(royaltyStatement.RoyaltyStatementUUID)
		if royaltyStatementExistingBytes == nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = "Royalty Statement does not exist!"
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementsOutput.FailureCount++
			continue
		}

		err = stub.PutState(royaltyStatement.RoyaltyStatementUUID, royaltyStatementsBytes)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
		}

		if royaltyStatementResponse.Success {
			royaltyStatementsOutput.SuccessCount++
		} else {
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementsOutput.FailureCount++
		}
	}

	royaltyStatementsOutput.RoyaltyStatement = royaltyStatementResponses

	objBytes, _ := objectToJSON(royaltyStatementsOutput)
	logger.Info("EXITING <", methodName, royaltyStatementsOutput)
	return shim.Success(objBytes)
}
