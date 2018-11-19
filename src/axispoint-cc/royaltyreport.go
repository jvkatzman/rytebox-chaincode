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

type RoyaltyReport struct {
	DocType         string `json:"docType"`
	RoyaltyReportID string `json:"royaltyReportID"`
}

//CreateRoyaltyReport : Record an royaltyReport
func CreateRoyaltyReport(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var Avalbytes []byte
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

	_ = jsonToObject([]byte(args[0]), royaltyReports)

	for _, royaltyReport := range *royaltyReports {
		royaltyReport.DocType = ROYALTYREPORT
		keys := []string{royaltyReport.RoyaltyReportID}
		royaltyReportResponse := RoyaltyReportResponse{}
		royaltyReportResponse.RoyaltyReport = royaltyReport
		royaltyReportResponse.Success = true

		//Check if royaltyReport already exists
		objBytes, err := QueryObject(stub, ROYALTYREPORT, keys)
		if err != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = err.Error()
		}
		if royaltyReportResponse.Success && objBytes != nil {
			royaltyReportResponse.Success = false
			royaltyReportResponse.Message = "Royalty Report with Royalty Report ID : " + royaltyReport.RoyaltyReportID + " already exists"
		} else {
			//Record royaltyReport on ledger
			Avalbytes, _ = objectToJSON(royaltyReport)
			err = UpdateObject(stub, ROYALTYREPORT, keys, Avalbytes)
			if err != nil {
				royaltyReportResponse.Success = false
				royaltyReportResponse.Message = err.Error()
			}
		}
		royaltyReportResponses = append(royaltyReportResponses, royaltyReportResponse)
	}
	objBytes, _ := objectToJSON(royaltyReportResponses)
	return shim.Success(objBytes)
}
