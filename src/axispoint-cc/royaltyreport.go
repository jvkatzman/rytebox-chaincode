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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// GetExploitationReportForQueryString : Get exploitation reports based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
var getExploitationReportForQueryString = getObjectByQueryFromLedger
var getRoyaltyReportsForQueryString = getObjectByQueryFromLedger

//AddRoyaltyReports : Add Royalty Reports to the ledger
func addRoyaltyReports(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addRoyaltyReports"
	logger.Info("ENTERING >", methodName, args)

	type RoyaltyReportResponse struct {
		RoyaltyReportUUID string `json:"royaltyReportUUID"`
		Message           string `json:"message"`
		Success           bool   `json:"success"`
	}

	type RoyaltyReportOutput struct {
		SuccessCount   int                     `json:"successCount"`
		FailureCount   int                     `json:"failureCount"`
		RoyaltyReports []RoyaltyReportResponse `json:"royaltyReports"`
	}

	if len(args) != 1 {
		return getErrorResponse("Missing arguments: Needed RoyaltyReport object to Create")
	}

	royaltyReportOutput := RoyaltyReportOutput{}
	royaltyReports := &[]RoyaltyReport{}
	royaltyReportResponses := []RoyaltyReportResponse{}

	err := jsonToObject([]byte(args[0]), royaltyReports)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Royalty Reports
	for _, royaltyReport := range *royaltyReports {
		royaltyReport.DocType = ROYALTYREPORT
		royaltyReportResponse := RoyaltyReportResponse{}
		royaltyReportResponse.RoyaltyReportUUID = royaltyReport.RoyaltyReportUUID
		royaltyReportResponse.Success = true

		exploitationReportUUID, err := getExploitationReportUUID(stub, royaltyReport)
		if err != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = err.Error()
			royaltyReportResponses = append(royaltyReportResponses, royaltyReportResponse)
			royaltyReportOutput.FailureCount++
			continue
		}

		royaltyReport.ExploitationReportUUID = exploitationReportUUID

		//Record royaltyReport on ledger
		royaltyReportBytes, err := objectToJSON(royaltyReport)
		if err != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = err.Error()
			royaltyReportResponses = append(royaltyReportResponses, royaltyReportResponse)
			royaltyReportOutput.FailureCount++
			continue
		}

		err = stub.PutState(royaltyReport.RoyaltyReportUUID, royaltyReportBytes)
		if err != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = err.Error()
		}

		if royaltyReportResponse.Success {
			royaltyReportOutput.SuccessCount++
		} else {
			royaltyReportOutput.FailureCount++
		}

		royaltyReportResponses = append(royaltyReportResponses, royaltyReportResponse)
	}

	royaltyReportOutput.RoyaltyReports = royaltyReportResponses

	objBytes, _ := objectToJSON(royaltyReportOutput)
	logger.Info("EXITING <", methodName, royaltyReportOutput)
	return shim.Success(objBytes)
}

//getExploitationReportUUID : Get the UUID of the exploitation report based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
func getExploitationReportUUID(stub shim.ChaincodeStubInterface, royaltyReport RoyaltyReport) (string, error) {
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
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\",\"target\":\"%s\",\"exploitationDate\":{\"$lte\":\"%s\"}}}", ROYALTYREPORT, targetCreator, exploitationDate)
	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getRoyaltyReportsForQueryString(stub, queryString) //getQueryResultInBytes(stub, queryString)
	if err != nil {
		errMsg := fmt.Sprintf("%s - Failed to get results for query: %s.  Error: %s", methodName, queryString, err.Error())
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}
	var resultRoyaltyReports []RoyaltyReport
	err = sliceToStruct(queryResult, &resultRoyaltyReports)
	if err != nil {
		errMsg := fmt.Sprintf("%s - Failed to convert query string result to struct: %s.  Error: %s", methodName, queryResult, err.Error())
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}
	queryResultBytes, err := json.Marshal(resultRoyaltyReports)
	if err != nil {
		errMsg := fmt.Sprintf("%s - Failed to convert royalty reports stuct to bytes.  Error: %s", methodName, err.Error())
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}
	logger.Infof("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}
