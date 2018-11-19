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
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//AddRoyaltyReports function contains business logic to insert new
// Royalty Reports to the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of exploitation report.
* @return   {pb.Response}    - peer Response
 */
func AddRoyaltyReports(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type RoyaltyReportResponse struct {
		RoyaltyReport RoyaltyReport `json:"royaltyReport"`
		Message       string        `json:"message"`
		Success       bool          `json:"success"`
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
		royaltyReportResponse.RoyaltyReport = royaltyReport
		royaltyReportResponse.Success = true

		//Record royaltyReport on ledger
		royaltyReportBytes, err := objectToJSON(royaltyReport)
		if err != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = err.Error()
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
	return shim.Success(objBytes)
}
