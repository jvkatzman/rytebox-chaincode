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

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// GetExploitationReportForQueryString function returns exploitation reports based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
var getExploitationReportForQueryString = getObjectByQueryFromLedger

//AddRoyaltyReports : Add Royalty Reports to the ledger
func addRoyaltyReports(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addRoyaltyReports"
	logger.Info("ENTERING >", methodName)

	type RoyaltyReportResponse struct {
		RoyaltyReportUUID string `json:"royaltyReportUUID"`
		Message           string `json:"message"`
		Success           bool   `json:"success"`
	}

	if len(args) != 1 {
		return getErrorResponse("Missing arguments: Needed RoyaltyReport object to Create")
	}

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
			continue
		}

		royaltyReport.ExploitationReportUUID = exploitationReportUUID

		//Record royaltyReport on ledger
		royaltyReportBytes, err := objectToJSON(royaltyReport)
		if err != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = err.Error()
			royaltyReportResponses = append(royaltyReportResponses, royaltyReportResponse)
			continue
		}

		err = stub.PutState(royaltyReport.RoyaltyReportUUID, royaltyReportBytes)
		if err != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = err.Error()
		}

		royaltyReportResponses = append(royaltyReportResponses, royaltyReportResponse)
	}

	objBytes, _ := objectToJSON(royaltyReportResponses)
	logger.Info("EXITING <", methodName)
	return shim.Success(objBytes)
}

//getExploitationReportUUID : Get the UUID of the exploitation report based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
func getExploitationReportUUID(stub shim.ChaincodeStubInterface, royaltyReport RoyaltyReport) (string, error) {
	var methodName = "getExploitationReportUUID"
	logger.Info("ENTERING >", methodName)

	exploitationReportUUID := ""
	queryString := "{\"selector\":{\"docType\":\"" + EXPLOITATIONREPORT + "\",\"source\": \"" + royaltyReport.Source + "\",\"isrc\": \"" + royaltyReport.Isrc + "\",\"exploitationDate\": \"" + royaltyReport.ExploitationDate + "\",\"territory\": \"" + royaltyReport.Territory + "\"}}"
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
		errorMessage := fmt.Sprintf("Cannot find Exploitation Report with Source: %s, ISRC: %s, Exploitation Date: %s, , Territory: %s", royaltyReport.Source, royaltyReport.Isrc, royaltyReport.ExploitationDate, royaltyReport.Territory)
		return exploitationReportUUID, errors.New(errorMessage)
	}

	exploitationReportUUID = exploitationReports[0].ExploitationReportUUID

	logger.Info("EXITING <", methodName, exploitationReportUUID)
	return exploitationReportUUID, nil
}
