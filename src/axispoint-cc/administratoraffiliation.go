package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var getAdministratorAffiliationsForQueryString = getObjectByQueryFromLedger

// AddAdministratorAffiliations function contains business logic to insert new
// Administrator Affiliations to the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of administrator affiliation.
* @return   {pb.Response}    - peer Response
 */
func addAdministratorAffiliations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addAdministratorAffiliations"
	logger.Info("ENTERING >", methodName, args)

	type AdministratorAffiliationResponse struct {
		AdministratorAffiliationUUID string `json:"administratorAffiliationUUID"`
		Message                      string `json:"message"`
		Success                      bool   `json:"success"`
	}

	type AdministratorAffiliationOutput struct {
		SuccessCount              int                                `json:"successCount"`
		FailureCount              int                                `json:"failureCount"`
		AdministratorAffiliations []AdministratorAffiliationResponse `json:"administratorAffiliations"`
	}

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Administrator Affiliation objects is required")
	}

	administratorAffiliationOutput := AdministratorAffiliationOutput{}
	administratorAffiliations := &[]AdministratorAffiliation{}
	administratorAffiliationResponses := []AdministratorAffiliationResponse{}

	//Unmarshal the args input to an array of administrator affiliation records
	err := jsonToObject([]byte(args[0]), administratorAffiliations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Exploitation Reports
	for _, administratorAffiliation := range *administratorAffiliations {
		administratorAffiliation.DocType = ADMINISTRATORAFFILIATION
		administratorAffiliationResponse := AdministratorAffiliationResponse{}
		administratorAffiliationResponse.AdministratorAffiliationUUID = administratorAffiliation.AdministratorAffiliationUUID
		administratorAffiliationResponse.Success = true

		//Record Exploitation Report on ledger
		administratorAffiliationBytes, err := objectToJSON(administratorAffiliation)
		if err != nil {
			administratorAffiliationResponse.Success = false
			administratorAffiliationResponse.Message = err.Error()
			administratorAffiliationResponses = append(administratorAffiliationResponses, administratorAffiliationResponse)
			administratorAffiliationOutput.FailureCount++
			continue
		}

		administratorAffiliationExistingBytes, err := stub.GetState(administratorAffiliation.AdministratorAffiliationUUID)
		if administratorAffiliationExistingBytes != nil {
			administratorAffiliationResponse.Success = false
			administratorAffiliationResponse.Message = "Administrator Affiliation already exists!"
			administratorAffiliationResponses = append(administratorAffiliationResponses, administratorAffiliationResponse)
			administratorAffiliationOutput.FailureCount++
			continue
		}

		err = stub.PutState(administratorAffiliation.AdministratorAffiliationUUID, administratorAffiliationBytes)
		if err != nil {
			administratorAffiliationResponse.Success = false
			administratorAffiliationResponse.Message = err.Error()
		}

		if administratorAffiliationResponse.Success {
			administratorAffiliationOutput.SuccessCount++
		} else {
			administratorAffiliationResponses = append(administratorAffiliationResponses, administratorAffiliationResponse)
			administratorAffiliationOutput.FailureCount++
		}
	}

	administratorAffiliationOutput.AdministratorAffiliations = administratorAffiliationResponses

	objBytes, _ := objectToJSON(administratorAffiliationOutput)
	logger.Info("EXITING <", methodName, administratorAffiliationOutput)
	return shim.Success(objBytes)
}

// updateAdministratorAffiliations function contains business logic to update
// Administrator Affiliations on the Ledger
/*
* @params   {Array} args
* @property {string} 0       - stringified JSON array of administrator affiliation.
* @return   {pb.Response}    - peer Response
 */
func updateAdministratorAffiliations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "updateAdministratorAffiliations"
	logger.Info("ENTERING >", methodName, args)

	type AdministratorAffiliationResponse struct {
		AdministratorAffiliationUUID string `json:"administratorAffiliationUUID"`
		Message                      string `json:"message"`
		Success                      bool   `json:"success"`
	}

	type AdministratorAffiliationOutput struct {
		SuccessCount              int                                `json:"successCount"`
		FailureCount              int                                `json:"failureCount"`
		AdministratorAffiliations []AdministratorAffiliationResponse `json:"administratorAffiliations"`
	}

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Administrator Affiliation objects is required")
	}

	administratorAffiliationOutput := AdministratorAffiliationOutput{}
	administratorAffiliations := &[]AdministratorAffiliation{}
	administratorAffiliationResponses := []AdministratorAffiliationResponse{}

	//Unmarshal the args input to an array of administrator affiliation records
	err := jsonToObject([]byte(args[0]), administratorAffiliations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// Iterate over Administrator Affiliations
	for _, administratorAffiliation := range *administratorAffiliations {
		administratorAffiliation.DocType = ADMINISTRATORAFFILIATION
		administratorAffiliationResponse := AdministratorAffiliationResponse{}
		administratorAffiliationResponse.AdministratorAffiliationUUID = administratorAffiliation.AdministratorAffiliationUUID
		administratorAffiliationResponse.Success = true

		//Record Administrator Affiliation on ledger
		administratorAffiliationBytes, err := objectToJSON(administratorAffiliation)
		if err != nil {
			administratorAffiliationResponse.Success = false
			administratorAffiliationResponse.Message = err.Error()
			administratorAffiliationResponses = append(administratorAffiliationResponses, administratorAffiliationResponse)
			administratorAffiliationOutput.FailureCount++
			continue
		}

		err = stub.PutState(administratorAffiliation.AdministratorAffiliationUUID, administratorAffiliationBytes)
		if err != nil {
			administratorAffiliationResponse.Success = false
			administratorAffiliationResponse.Message = err.Error()
		}

		if administratorAffiliationResponse.Success {
			administratorAffiliationOutput.SuccessCount++
		} else {
			administratorAffiliationResponses = append(administratorAffiliationResponses, administratorAffiliationResponse)
			administratorAffiliationOutput.FailureCount++
		}
	}

	administratorAffiliationOutput.AdministratorAffiliations = administratorAffiliationResponses

	objBytes, _ := objectToJSON(administratorAffiliationOutput)
	logger.Info("EXITING <", methodName, administratorAffiliationOutput)
	return shim.Success(objBytes)
}

//getAdministratorAffiliations: get administrator affiliations
func getAdministratorAffiliations(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getAdministratorAffiliations"
	logger.Info("ENTERING >", methodName, args)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", ADMINISTRATORAFFILIATION)
	if len(args) == 1 {
		queryString = args[0]
	}

	logger.Info("%s - executing rich query : %s.", methodName, queryString)

	queryResult, err := getAdministratorAffiliationsForQueryString(stub, queryString) //getQueryResultInBytes(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	var resultAdministratorAffiliations []AdministratorAffiliation
	err = sliceToStruct(queryResult, &resultAdministratorAffiliations)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	queryResultBytes, err := objectToJSON(resultAdministratorAffiliations)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Info("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}
