package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var getHolderRepresentationsForQueryString = getObjectByQueryFromLedger

// AddHolderRepresentations function contains business logic to insert new
// Holder Representations to the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of holder representation.
* @return   {pb.Response}    - peer Response
 */
func addHolderRepresentations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addHolderRepresentations"
	logger.Info("ENTERING >", methodName, args)

	type HolderRepresentationResponse struct {
		HolderRepresentationUUID string `json:"holderRepresentationUUID"`
		Message                  string `json:"message"`
		Success                  bool   `json:"success"`
	}

	type HolderRepresentationOutput struct {
		SuccessCount          int                            `json:"successCount"`
		FailureCount          int                            `json:"failureCount"`
		HolderRepresentations []HolderRepresentationResponse `json:"holderRepresentations"`
	}

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Holder Representation objects is required")
	}

	holderRepresentationOutput := HolderRepresentationOutput{}
	holderRepresentations := &[]HolderRepresentation{}
	holderRepresentationResponses := []HolderRepresentationResponse{}

	//Unmarshal the args input to an array of holder representation records
	err := jsonToObject([]byte(args[0]), holderRepresentations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Exploitation Reports
	for _, holderRepresentation := range *holderRepresentations {
		holderRepresentation.DocType = HOLDERREPRESENTATION
		holderRepresentationResponse := HolderRepresentationResponse{}
		holderRepresentationResponse.HolderRepresentationUUID = holderRepresentation.HolderRepresentationUUID
		holderRepresentationResponse.Success = true

		//Record Exploitation Report on ledger
		holderRepresentationBytes, err := objectToJSON(holderRepresentation)
		if err != nil {
			holderRepresentationResponse.Success = false
			holderRepresentationResponse.Message = err.Error()
			holderRepresentationResponses = append(holderRepresentationResponses, holderRepresentationResponse)
			holderRepresentationOutput.FailureCount++
			continue
		}

		holderRepresentationExistingBytes, err := stub.GetState(holderRepresentation.HolderRepresentationUUID)
		if holderRepresentationExistingBytes != nil {
			holderRepresentationResponse.Success = false
			holderRepresentationResponse.Message = "Holder Representation already exists!"
			holderRepresentationResponses = append(holderRepresentationResponses, holderRepresentationResponse)
			holderRepresentationOutput.FailureCount++
			continue
		}

		err = stub.PutState(holderRepresentation.HolderRepresentationUUID, holderRepresentationBytes)
		if err != nil {
			holderRepresentationResponse.Success = false
			holderRepresentationResponse.Message = err.Error()
		}

		if holderRepresentationResponse.Success {
			holderRepresentationOutput.SuccessCount++
		} else {
			holderRepresentationResponses = append(holderRepresentationResponses, holderRepresentationResponse)
			holderRepresentationOutput.FailureCount++
		}
	}

	holderRepresentationOutput.HolderRepresentations = holderRepresentationResponses

	objBytes, _ := objectToJSON(holderRepresentationOutput)
	logger.Info("EXITING <", methodName, holderRepresentationOutput)
	return shim.Success(objBytes)
}

// updateHolderRepresentations function contains business logic to update
// Holder Representations on the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of holder representation.
* @return   {pb.Response}    - peer Response
 */
func updateHolderRepresentations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "updateHolderRepresentations"
	logger.Info("ENTERING >", methodName, args)

	type HolderRepresentationResponse struct {
		HolderRepresentationUUID string `json:"holderRepresentationUUID"`
		Message                  string `json:"message"`
		Success                  bool   `json:"success"`
	}

	type HolderRepresentationOutput struct {
		SuccessCount          int                            `json:"successCount"`
		FailureCount          int                            `json:"failureCount"`
		HolderRepresentations []HolderRepresentationResponse `json:"holderRepresentations"`
	}

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Holder Representation objects is required")
	}

	holderRepresentationOutput := HolderRepresentationOutput{}
	holderRepresentations := &[]HolderRepresentation{}
	holderRepresentationResponses := []HolderRepresentationResponse{}

	//Unmarshal the args input to an array of holder representation records
	err := jsonToObject([]byte(args[0]), holderRepresentations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Exploitation Reports
	for _, holderRepresentation := range *holderRepresentations {
		holderRepresentation.DocType = HOLDERREPRESENTATION
		holderRepresentationResponse := HolderRepresentationResponse{}
		holderRepresentationResponse.HolderRepresentationUUID = holderRepresentation.HolderRepresentationUUID
		holderRepresentationResponse.Success = true

		//Record Exploitation Report on ledger
		holderRepresentationBytes, err := objectToJSON(holderRepresentation)
		if err != nil {
			holderRepresentationResponse.Success = false
			holderRepresentationResponse.Message = err.Error()
			holderRepresentationResponses = append(holderRepresentationResponses, holderRepresentationResponse)
			holderRepresentationOutput.FailureCount++
			continue
		}

		err = stub.PutState(holderRepresentation.HolderRepresentationUUID, holderRepresentationBytes)
		if err != nil {
			holderRepresentationResponse.Success = false
			holderRepresentationResponse.Message = err.Error()
		}

		if holderRepresentationResponse.Success {
			holderRepresentationOutput.SuccessCount++
		} else {
			holderRepresentationResponses = append(holderRepresentationResponses, holderRepresentationResponse)
			holderRepresentationOutput.FailureCount++
		}
	}

	holderRepresentationOutput.HolderRepresentations = holderRepresentationResponses

	objBytes, _ := objectToJSON(holderRepresentationOutput)
	logger.Info("EXITING <", methodName, holderRepresentationOutput)
	return shim.Success(objBytes)
}

//getHolderRepresentations: get holder representations
func getHolderRepresentationByUUID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getHolderRepresentationByUUID"
	logger.Info("ENTERING >", methodName, args)

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Holder Representation UUID is missing")
	}

	holderRepresentationBytes, err := stub.GetState(args[0])
	if err != nil {
		return getErrorResponse(err.Error())
	}
	if holderRepresentationBytes == nil {
		return getErrorResponse(fmt.Sprintf("Holder Representation with UUID: %s does not exist", args[0]))
	}

	//return bytes as result
	return shim.Success(holderRepresentationBytes)
}

//getHolderRepresentations: get holder representations
func getHolderRepresentations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getHolderRepresentations"
	logger.Info("ENTERING >", methodName, args)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", HOLDERREPRESENTATION)
	if len(args) > 1 {
		queryString = args[0]
	}

	logger.Info("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getHolderRepresentationsForQueryString(stub, queryString) //getQueryResultInBytes(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	var resultHolderRepresentations []HolderRepresentation
	err = sliceToStruct(queryResult, &resultHolderRepresentations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	queryResultBytes, err := objectToJSON(resultHolderRepresentations)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Info("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}
