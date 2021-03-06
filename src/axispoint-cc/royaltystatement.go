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

// RoyaltyStatementResponse : defines response data from blockchain request
type RoyaltyStatementResponse struct {
	RoyaltyStatementUUID string `json:"royaltyStatementUUID"`
	Message              string `json:"message"`
	Success              bool   `json:"success"`
}

// RoyaltyStatementOutput : defines accumulated output of blockchain requests
type RoyaltyStatementOutput struct {
	SuccessCount      int                        `json:"successCount"`
	FailureCount      int                        `json:"failureCount"`
	RoyaltyStatements []RoyaltyStatementResponse `json:"royaltyStatements"`
}

// getExploitationReportForQueryString : Get exploitation reports based on Song Title, Song Writer, ISRC, Exploitation Date and Territory
var getExploitationReportForQueryString = getObjectByQueryFromLedger

// getRoyaltyStatementsForQueryString : Get royalty statements based on rich query selectoe
var getRoyaltyStatementsForQueryString = getObjectByQueryFromLedger

/* addRoyaltyStatements function contains business logic to insert new
Royalty Statements to the Ledger
* @params   {Array} args
* @property {string} 0       - stringified JSON array of royalty statement.
* @return   {pb.Response}    - peer Response
*/
// refactor the following : create 2 separete methods
//1.  Save royaltyStatement As is
//2.  Save + generate Event.
func addRoyaltyStatements(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addRoyaltyStatements"
	logger.Info("ENTERING >", methodName, args)

	//if this function is called with morethan 1 royalty statements
	//only write if the IPI or the orgs are same for 2 or  more royalty statements
	//otherwise the chaincode should not continue.
	royaltyStatementsEventPayloadBytes := []byte{}

	if len(args) != 1 {
		return getErrorResponse("Missing arguments: Needed RoyaltyStatement object to Create")
	}

	royaltyStatementOutput := RoyaltyStatementOutput{}
	royaltyStatements := &[]RoyaltyStatement{}
	royaltyStatementResponses := []RoyaltyStatementResponse{}

	// Unmarshal the args input to an array of royalty statement records
	err := jsonToObject([]byte(args[0]), royaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// iterate over royalty statements
	for _, royaltyStatement := range *royaltyStatements {
		royaltyStatement.DocType = ROYALTYSTATEMENT
		royaltyStatementResponse := RoyaltyStatementResponse{}
		royaltyStatementResponse.RoyaltyStatementUUID = royaltyStatement.RoyaltyStatementUUID
		royaltyStatementResponse.Success = true

		// check if royalty statement already exists
		royaltyStatementExistingBytes, err := stub.GetState(royaltyStatement.RoyaltyStatementUUID)
		if royaltyStatementExistingBytes != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = "Royalty Statement already exists!"
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
			continue
		}

		// convert royalty statement to bytes
		royaltyStatementBytes, err := objectToJSON(royaltyStatement)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
			continue
		}

		// add royalty statement to the ledger
		err = stub.PutState(royaltyStatement.RoyaltyStatementUUID, royaltyStatementBytes)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
		} else {
			//assumption: this method should be called with a single royalty statement for now.
			//TODO: multiple royalty statements if ORG or IPI is the same.
			payloadBytes, err := getRoyaltyStatementsEventPayload(stub, royaltyStatement)
			if err != nil {
				return getErrorResponse(fmt.Sprintf("%s - Failed to construct '%s' payload.  Error: %s", methodName, EventRoyaltyStatementCreation, err.Error()))
			}
			royaltyStatementsEventPayloadBytes = append(royaltyStatementsEventPayloadBytes, payloadBytes...)
		}

		if royaltyStatementResponse.Success {
			royaltyStatementOutput.SuccessCount++
		} else {
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
		}
	}

	royaltyStatementOutput.RoyaltyStatements = royaltyStatementResponses

	objBytes, _ := objectToJSON(royaltyStatementOutput)

	//fire an event for Ownership report only
	if len(*royaltyStatements) == 1 && (*royaltyStatements)[0].RightType == OWNERSHIP {
		logger.Infof("%s - firing event '%s'.", methodName, EventRoyaltyStatementCreation)
		err = stub.SetEvent(EventRoyaltyStatementCreation, royaltyStatementsEventPayloadBytes)
		if err != nil {
			return getErrorResponse(fmt.Sprintf("%s - Failed to set event '%s' with payload '%s'.  Error: %s", methodName, EventRoyaltyStatementCreation, royaltyStatementsEventPayloadBytes, err.Error()))
		}
	}

	logger.Info("EXITING <", methodName, royaltyStatementOutput)
	return shim.Success(objBytes)
}

//addRoyaltyStatementsAndEvent - save the royalty statement and fire an event
func addRoyaltyStatementAndEvent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addRoyaltyStatementAndEvent"
	logger.Info("ENTERING >", methodName, args)

	//if this function is called with morethan 1 royalty statements
	//only write if the IPI or the orgs are same for 2 or  more royalty statements
	//otherwise the chaincode should not continue.
	royaltyStatementsEventPayloadBytes := []byte{}
	isFinalRoyaltyStatement := false

	if len(args) != 1 {
		return getErrorResponse("Missing arguments: Needed RoyaltyStatement object to Create")
	}

	royaltyStatementOutput := RoyaltyStatementOutput{}
	royaltyStatements := &[]RoyaltyStatement{}
	royaltyStatementResponses := []RoyaltyStatementResponse{}

	// Unmarshal the args input to an array of royalty statement records
	err := jsonToObject([]byte(args[0]), royaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// iterate over royalty statements
	for _, royaltyStatement := range *royaltyStatements {
		royaltyStatement.DocType = ROYALTYSTATEMENT
		royaltyStatementResponse := RoyaltyStatementResponse{}
		royaltyStatementResponse.RoyaltyStatementUUID = royaltyStatement.RoyaltyStatementUUID
		royaltyStatementResponse.Success = true

		// check if royalty statement already exists
		royaltyStatementExistingBytes, err := stub.GetState(royaltyStatement.RoyaltyStatementUUID)
		if royaltyStatementExistingBytes != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = "Royalty Statement already exists!"
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
			continue
		}

		// convert royalty statement to bytes
		royaltyStatementBytes, err := objectToJSON(royaltyStatement)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
			continue
		}

		if royaltyStatement.CollectionRight == 0 && royaltyStatement.CollectionRightPercent == 0 {
			isFinalRoyaltyStatement = true
			logger.Infof("%s - final royalty statement received with uuid : %s", methodName, royaltyStatement.RoyaltyStatementUUID)
		}

		// add royalty statement to the ledger
		err = stub.PutState(royaltyStatement.RoyaltyStatementUUID, royaltyStatementBytes)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
		} else if isFinalRoyaltyStatement == false {
			//assumption: this method should be called with a single royalty statement for now.
			//TODO: multiple royalty statements if ORG or IPI is the same.
			payloadBytes, err := getRoyaltyStatementsEventPayload(stub, royaltyStatement)
			if err != nil {
				return getErrorResponse(fmt.Sprintf("%s - Failed to construct '%s' payload.  Error: %s", methodName, EventRoyaltyStatementCreation, err.Error()))
			}
			royaltyStatementsEventPayloadBytes = append(royaltyStatementsEventPayloadBytes, payloadBytes...)
		}

		if royaltyStatementResponse.Success {
			royaltyStatementOutput.SuccessCount++
		} else {
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementOutput.FailureCount++
		}
	}

	royaltyStatementOutput.RoyaltyStatements = royaltyStatementResponses

	objBytes, _ := objectToJSON(royaltyStatementOutput)

	//fire an event for any royalty statements as long as its not the last one.
	if isFinalRoyaltyStatement == false {
		err = stub.SetEvent(EventRoyaltyStatementCreation, royaltyStatementsEventPayloadBytes)
		if err != nil {
			return getErrorResponse(fmt.Sprintf("%s - Failed to set event '%s' with payload '%s'.  Error: %s", methodName, EventRoyaltyStatementCreation, royaltyStatementsEventPayloadBytes, err.Error()))
		}
	} else {
		logger.Infof("%s - event '%s' not fired due to final royalty report.", methodName, EventRoyaltyStatementCreation)
	}

	logger.Info("EXITING <", methodName, royaltyStatementOutput)
	return shim.Success(objBytes)
}

/* getRoyaltyStatements function contains business logic to get
Royalty Statements based on the rich query selector
* @params   {Array} args
* @property {string} 0       - rich query selector.
* @return   {pb.Response}    - peer Response
*/
func getRoyaltyStatements(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getRoyaltyStatements"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", ROYALTYSTATEMENT)
	if len(args) == 1 {
		queryString = args[0]
	}

	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	// get royalty statements based on the rich query selector
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

/* getRoyaltyStatementsByUUIDs function contains business logic to get
Royalty Statements based on UUIDs
* @params   {Array} args	 - list of UUIDs
* @return   {pb.Response}    - peer Response
*/
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

/* updateRoyaltyStatements function contains business logic to update
Royalty Statements on the Ledger
* @params   {Array} args
* @property {string} 0       - stringified JSON array of royalty statement.
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

	// Unmarshal the args input to an array of royalty statement records
	err := jsonToObject([]byte(args[0]), royaltyStatements)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// iterate over royalty statements
	for _, royaltyStatement := range *royaltyStatements {
		royaltyStatement.DocType = ROYALTYSTATEMENT
		royaltyStatementResponse := RoyaltyStatementResponse{}
		royaltyStatementResponse.RoyaltyStatementUUID = royaltyStatement.RoyaltyStatementUUID
		royaltyStatementResponse.Success = true

		// convert royalty statement to bytes
		royaltyStatementsBytes, err := objectToJSON(royaltyStatement)
		if err != nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = err.Error()
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementsOutput.FailureCount++
			continue
		}

		// check if royalty statement with the UUID exists on the ledger.
		royaltyStatementExistingBytes, err := stub.GetState(royaltyStatement.RoyaltyStatementUUID)
		if royaltyStatementExistingBytes == nil {
			royaltyStatementResponse.Success = false
			royaltyStatementResponse.Message = "Royalty Statement does not exist!"
			royaltyStatementResponses = append(royaltyStatementResponses, royaltyStatementResponse)
			royaltyStatementsOutput.FailureCount++
			continue
		}

		// update royalty statement on the ledger
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

	royaltyStatementsOutput.RoyaltyStatements = royaltyStatementResponses

	objBytes, _ := objectToJSON(royaltyStatementsOutput)
	logger.Info("EXITING <", methodName, royaltyStatementsOutput)
	return shim.Success(objBytes)
}

/* updateRoyaltyStatements function contains business logic to update
Royalty Statements on the Ledger
* @params   {Array} args
* @property {string} 0       - stringified JSON array of royalty statement.
* @return   {pb.Response}    - peer Response
*/
func getExploitationReportUUID(stub shim.ChaincodeStubInterface, royaltyStatement RoyaltyStatement) (string, error) {
	var methodName = "getExploitationReportUUID"
	logger.Info("ENTERING >", methodName)

	exploitationReportUUID := ""
	queryString := "{\"selector\":{\"docType\":\"" + EXPLOITATIONREPORT + "\",\"source\": \"" + royaltyStatement.Source + "\",\"isrc\": \"" + royaltyStatement.Isrc + "\",\"exploitationDate\": \"" + royaltyStatement.ExploitationDate + "\",\"territory\": \"" + royaltyStatement.Territory + "\",\"usageType\": \"" + royaltyStatement.UsageType + "\"}}"
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
		errorMessage := fmt.Sprintf("Cannot find Exploitation Report with Source: %s, ISRC: %s, Exploitation Date: %s, Territory: %s, Usage Type: %s", royaltyStatement.Source, royaltyStatement.Isrc, royaltyStatement.ExploitationDate, royaltyStatement.Territory, royaltyStatement.UsageType)
		return exploitationReportUUID, errors.New(errorMessage)
	}

	exploitationReportUUID = exploitationReports[0].ExploitationReportUUID

	logger.Info("EXITING <", methodName, exploitationReportUUID)
	return exploitationReportUUID, nil
}

func getRoyaltyStatementsEventPayload(stub shim.ChaincodeStubInterface, royaltyStatement RoyaltyStatement /*royaltyStatements *[]RoyaltyStatement*/) ([]byte, error) {
	methodName := "getRoyaltyStatementsEventPayload"

	objRoyaltyStatementEventPayload := RoyaltyStatementCreationEventPayload{}
	objRoyaltyStatementEventPayload.RoyaltyStatementUUID = royaltyStatement.RoyaltyStatementUUID
	objRoyaltyStatementEventPayload.Type = COLLECTION
	logger.Infof("%s - setting 'type' for event payload to '%s'.", methodName, objRoyaltyStatementEventPayload.Type)
	/*if len(royaltyStatement.RightHolder) > 0 && len(royaltyStatement.Administrator) > 0 {
		logger.Infof("%s - found a valid royalty statement right holder '%s' and administrator '%s'.", methodName, royaltyStatement.RightHolder, royaltyStatement.Administrator)
		objRoyaltyStatementEventPayload.TargetIPI = royaltyStatement.Administrator
		objRoyaltyStatementEventPayload.IsDSP = false
	} else if len(royaltyStatement.Collector) > 0 && len(royaltyStatement.Administrator) > 0 {
		logger.Infof("%s - found a valid royalty statement collector '%s' and administrator '%s'.", methodName, royaltyStatement.Collector, royaltyStatement.Administrator)
		objRoyaltyStatementEventPayload.TargetIPI = royaltyStatement.Collector
		objRoyaltyStatementEventPayload.IsDSP = false

	}*/
	if len(royaltyStatement.Collector) > 0 && len(royaltyStatement.Administrator) > 0 {
		logger.Infof("%s - found a valid royalty statement collector '%s' and administrator '%s'.", methodName, royaltyStatement.Collector, royaltyStatement.Administrator)
		objRoyaltyStatementEventPayload.TargetIPI = royaltyStatement.Collector
		objRoyaltyStatementEventPayload.IsDSP = false

	} else if len(royaltyStatement.RightHolder) > 0 && len(royaltyStatement.Administrator) > 0 {
		logger.Infof("%s - found a valid royalty statement right holder '%s' and administrator '%s'.", methodName, royaltyStatement.RightHolder, royaltyStatement.Administrator)
		objRoyaltyStatementEventPayload.TargetIPI = royaltyStatement.Administrator
		objRoyaltyStatementEventPayload.IsDSP = false
	} else if len(royaltyStatement.Collector) > 0 && len(royaltyStatement.Administrator) == 0 {
		logger.Infof("%s - found a valid royalty statement collector '%s' and an invalid administrator.  Using source '%s'.", methodName, royaltyStatement.Collector, royaltyStatement.Source)
		objRoyaltyStatementEventPayload.TargetIPI = royaltyStatement.Source
		objRoyaltyStatementEventPayload.IsDSP = true
	} else if royaltyStatement.RightType == OWNERSHIP {
		objRoyaltyStatementEventPayload.TargetIPI = royaltyStatement.RightHolder
		objRoyaltyStatementEventPayload.IsDSP = false
		objRoyaltyStatementEventPayload.Type = OWNERSHIP

	} else {
		//error condition
		message := fmt.Sprintf("%s - incorrect condition found when determining the target IPI and dsp status.  Operation cannot continue", methodName)
		logger.Error(message)
		return nil, errors.New(message)
	}
	//get the org from the mapping stored on the chain
	//ipiToOrgBytes, err := getAssetByUUID(stub, []string{royaltyStatement.RightHolder}]).//stub.GetState(royaltyStatement.RightHolder)
	response := getAssetByUUID(stub, []string{royaltyStatement.RightHolder})
	if response.GetStatus() != shim.OK {
		//if err != nil {
		//message := fmt.Sprintf("%s - Failed to get org for IPI '%s'.  Error: %s", methodName, royaltyStatement.RightHolder, err.Error())
		message := fmt.Sprintf("%s - Failed to get org for IPI '%s'.  Error: %s", methodName, royaltyStatement.RightHolder, response.GetMessage())
		logger.Error(message)
		return nil, errors.New(message)
	}
	ipiToOrgBytes := response.GetPayload()

	objRoyaltyStatementEventPayload.TargetOrg = string(ipiToOrgBytes) //`{"org":"org2"}`
	logger.Infof("%s - setting 'target org' for event payload to '%s'.", methodName, objRoyaltyStatementEventPayload.TargetOrg)

	objResultBytes, err := objectToJSON(objRoyaltyStatementEventPayload)
	if err != nil {
		message := fmt.Sprintf("%s - Failed get payload in bytes for IPI to org mapping.  Error: %s", methodName, err.Error())
		logger.Error(message)
		return nil, errors.New(message)
	}
	logger.Infof("%s - event payload being returned: %s", methodName, string(objResultBytes))
	return objResultBytes, nil
}
