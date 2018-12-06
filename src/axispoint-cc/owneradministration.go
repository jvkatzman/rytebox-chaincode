package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var getOwnerAdministrationsForQueryString = getObjectByQueryFromLedger

// AddOwnerAdministrations function contains business logic to insert new
// Owner Administrations to the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of owner administration.
* @return   {pb.Response}    - peer Response
 */
func addOwnerAdministrations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addOwnerAdministrations"
	logger.Info("ENTERING >", methodName, args)

	type OwnerAdministrationResponse struct {
		OwnerAdministrationUUID string `json:"ownerAdministrationUUID"`
		Message                 string `json:"message"`
		Success                 bool   `json:"success"`
	}

	type OwnerAdministrationOutput struct {
		SuccessCount         int                           `json:"successCount"`
		FailureCount         int                           `json:"failureCount"`
		OwnerAdministrations []OwnerAdministrationResponse `json:"ownerAdministrations"`
	}

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Owner Administration objects is required")
	}

	ownerAdministrationOutput := OwnerAdministrationOutput{}
	ownerAdministrations := &[]OwnerAdministration{}
	ownerAdministrationResponses := []OwnerAdministrationResponse{}

	//Unmarshal the args input to an array of owner administration records
	err := jsonToObject([]byte(args[0]), ownerAdministrations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Exploitation Reports
	for _, ownerAdministration := range *ownerAdministrations {
		ownerAdministration.DocType = HOLDERREPRESENTATION
		ownerAdministrationResponse := OwnerAdministrationResponse{}
		ownerAdministrationResponse.OwnerAdministrationUUID = ownerAdministration.OwnerAdministrationUUID
		ownerAdministrationResponse.Success = true

		//Record Exploitation Report on ledger
		ownerAdministrationBytes, err := objectToJSON(ownerAdministration)
		if err != nil {
			ownerAdministrationResponse.Success = false
			ownerAdministrationResponse.Message = err.Error()
			ownerAdministrationResponses = append(ownerAdministrationResponses, ownerAdministrationResponse)
			ownerAdministrationOutput.FailureCount++
			continue
		}

		ownerAdministrationExistingBytes, err := stub.GetState(ownerAdministration.OwnerAdministrationUUID)
		if ownerAdministrationExistingBytes != nil {
			ownerAdministrationResponse.Success = false
			ownerAdministrationResponse.Message = "Owner Administration already exists!"
			ownerAdministrationResponses = append(ownerAdministrationResponses, ownerAdministrationResponse)
			ownerAdministrationOutput.FailureCount++
			continue
		}

		err = stub.PutState(ownerAdministration.OwnerAdministrationUUID, ownerAdministrationBytes)
		if err != nil {
			ownerAdministrationResponse.Success = false
			ownerAdministrationResponse.Message = err.Error()
		}

		if ownerAdministrationResponse.Success {
			ownerAdministrationOutput.SuccessCount++
		} else {
			ownerAdministrationResponses = append(ownerAdministrationResponses, ownerAdministrationResponse)
			ownerAdministrationOutput.FailureCount++
		}
	}

	ownerAdministrationOutput.OwnerAdministrations = ownerAdministrationResponses

	objBytes, _ := objectToJSON(ownerAdministrationOutput)
	logger.Info("EXITING <", methodName, ownerAdministrationOutput)
	return shim.Success(objBytes)
}

// updateOwnerAdministrations function contains business logic to update
// Owner Administrations on the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of owner administration.
* @return   {pb.Response}    - peer Response
 */
func updateOwnerAdministrations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "updateOwnerAdministrations"
	logger.Info("ENTERING >", methodName, args)

	type OwnerAdministrationResponse struct {
		OwnerAdministrationUUID string `json:"ownerAdministrationUUID"`
		Message                 string `json:"message"`
		Success                 bool   `json:"success"`
	}

	type OwnerAdministrationOutput struct {
		SuccessCount         int                           `json:"successCount"`
		FailureCount         int                           `json:"failureCount"`
		OwnerAdministrations []OwnerAdministrationResponse `json:"ownerAdministrations"`
	}

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Owner Administration objects is required")
	}

	ownerAdministrationOutput := OwnerAdministrationOutput{}
	ownerAdministrations := &[]OwnerAdministration{}
	ownerAdministrationResponses := []OwnerAdministrationResponse{}

	//Unmarshal the args input to an array of owner administration records
	err := jsonToObject([]byte(args[0]), ownerAdministrations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Exploitation Reports
	for _, ownerAdministration := range *ownerAdministrations {
		ownerAdministration.DocType = HOLDERREPRESENTATION
		ownerAdministrationResponse := OwnerAdministrationResponse{}
		ownerAdministrationResponse.OwnerAdministrationUUID = ownerAdministration.OwnerAdministrationUUID
		ownerAdministrationResponse.Success = true

		//Record Exploitation Report on ledger
		ownerAdministrationBytes, err := objectToJSON(ownerAdministration)
		if err != nil {
			ownerAdministrationResponse.Success = false
			ownerAdministrationResponse.Message = err.Error()
			ownerAdministrationResponses = append(ownerAdministrationResponses, ownerAdministrationResponse)
			ownerAdministrationOutput.FailureCount++
			continue
		}

		err = stub.PutState(ownerAdministration.OwnerAdministrationUUID, ownerAdministrationBytes)
		if err != nil {
			ownerAdministrationResponse.Success = false
			ownerAdministrationResponse.Message = err.Error()
		}

		if ownerAdministrationResponse.Success {
			ownerAdministrationOutput.SuccessCount++
		} else {
			ownerAdministrationResponses = append(ownerAdministrationResponses, ownerAdministrationResponse)
			ownerAdministrationOutput.FailureCount++
		}
	}

	ownerAdministrationOutput.OwnerAdministrations = ownerAdministrationResponses

	objBytes, _ := objectToJSON(ownerAdministrationOutput)
	logger.Info("EXITING <", methodName, ownerAdministrationOutput)
	return shim.Success(objBytes)
}

//getOwnerAdministrations: get owner administrations
func getOwnerAdministrations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getOwnerAdministrations"
	logger.Info("ENTERING >", methodName, args)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", HOLDERREPRESENTATION)
	if len(args) == 1 {
		queryString = args[0]
	}

	logger.Info("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getOwnerAdministrationsForQueryString(stub, queryString) //getQueryResultInBytes(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	var resultOwnerAdministrations []OwnerAdministration
	err = sliceToStruct(queryResult, &resultOwnerAdministrations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	queryResultBytes, err := objectToJSON(resultOwnerAdministrations)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Info("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}
