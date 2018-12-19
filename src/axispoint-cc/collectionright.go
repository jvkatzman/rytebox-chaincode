package main

import (
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

var getCollectionRightsForQueryString = getObjectByQueryFromLedger

// CollectionRightsResponse : defines response data from blockchain request
type CollectionRightsResponse struct {
	CollectionRightUUID string `json:"collectionRightUUID"`
	Message             string `json:"message"`
	Success             bool   `json:"success"`
}

// CollectionRightsOutput : defines accumulated output of blockchain requests
type CollectionRightsOutput struct {
	SuccessCount              int                        `json:"successCount"`
	FailureCount              int                        `json:"failureCount"`
	CollectionRightsResponses []CollectionRightsResponse `json:"collectionRightsResponses"`
}

/* addCollectionRights function contains business logic to insert new
Collection Reports to the Ledger
* @params   {Array} args
* @property {string} 0       - stringified JSON array of royalty statement.
* @return   {pb.Response}    - peer Response
*/
func addCollectionRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addCollectionRights"
	logger.Info("ENTERING >", methodName, args)

	if len(args) != 1 {
		return getErrorResponse("Missing arguments: Needed Collection Reports object to Create")
	}

	collectionRightsOutput := CollectionRightsOutput{}
	collectionRights := &[]CollectionRight{}
	collectionRightsResponses := []CollectionRightsResponse{}

	// Unmarshal the args input to an array of royalty statement records
	err := jsonToObject([]byte(args[0]), collectionRights)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// iterate over royalty statements
	for _, collectionRight := range *collectionRights {
		collectionRight.DocType = COLLECTIONRIGHTSREPORT
		collectionRightsResponse := CollectionRightsResponse{}
		collectionRightsResponse.CollectionRightUUID = collectionRight.CollectionRightUUID
		collectionRightsResponse.Success = true

		// check if royalty statement already exists
		collectionRightExistingBytes, err := stub.GetState(collectionRight.CollectionRightUUID)
		if collectionRightExistingBytes != nil {
			collectionRightsResponse.Success = false
			collectionRightsResponse.Message = "Collection Right already exists!"
			collectionRightsResponses = append(collectionRightsResponses, collectionRightsResponse)
			collectionRightsOutput.FailureCount++
			continue
		}

		// convert royalty statement to bytes
		collectionRightBytes, err := objectToJSON(collectionRight)
		if err != nil {
			collectionRightsResponse.Success = false
			collectionRightsResponse.Message = err.Error()
			collectionRightsResponses = append(collectionRightsResponses, collectionRightsResponse)
			collectionRightsOutput.FailureCount++
			continue
		}

		// add royalty statement to the ledger
		err = stub.PutState(collectionRight.CollectionRightUUID, collectionRightBytes)
		if err != nil {
			collectionRightsResponse.Success = false
			collectionRightsResponse.Message = err.Error()
		}

		if collectionRightsResponse.Success {
			collectionRightsOutput.SuccessCount++
		} else {
			collectionRightsResponses = append(collectionRightsResponses, collectionRightsResponse)
			collectionRightsOutput.FailureCount++
		}
	}

	collectionRightsOutput.CollectionRightsResponses = collectionRightsResponses

	objBytes, _ := objectToJSON(collectionRightsOutput)
	logger.Info("EXITING <", methodName, collectionRightsOutput)
	return shim.Success(objBytes)
}

/* getCollectionRights function contains business logic to get
Royalty Statements based on the rich query selector
* @params   {Array} args
* @property {string} 0       - rich query selector.
* @return   {pb.Response}    - peer Response
*/
func getCollectionRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getCollectionRights"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", COLLECTIONRIGHTSREPORT)
	if len(args) == 1 {
		queryString = args[0]
	}

	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	// get royalty statements based on the rich query selector
	queryResult, err := getCollectionRightsForQueryString(stub, queryString)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	var resultCollectionRights []CollectionRight
	err = sliceToStruct(queryResult, &resultCollectionRights)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	queryResultBytes, err := objectToJSON(resultCollectionRights)
	if err != nil {
		return getErrorResponse(err.Error())
	}
	logger.Infof("result(s) received from couch db: %s", string(queryResultBytes))

	//return bytes as result
	return shim.Success(queryResultBytes)
}

/* updateCollectionRights function contains business logic to update
Royalty Statements on the Ledger
* @params   {Array} args
* @property {string} 0       - stringified JSON array of royalty statement.
* @return   {pb.Response}    - peer Response
*/
func updateCollectionRights(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "updateCollectionRights"
	logger.Info("ENTERING >", methodName, args)

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: Array of Collection Right objects is required")
	}

	collectionRightsOutput := CollectionRightsOutput{}
	collectionRights := &[]CollectionRight{}
	collectionRightsResponses := []CollectionRightsResponse{}

	// Unmarshal the args input to an array of collection rights
	err := jsonToObject([]byte(args[0]), collectionRights)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	// iterate over royalty statements
	for _, collectionRight := range *collectionRights {
		collectionRight.DocType = COLLECTIONRIGHTSREPORT
		collectionRightsResponse := CollectionRightsResponse{}
		collectionRightsResponse.CollectionRightUUID = collectionRight.CollectionRightUUID
		collectionRightsResponse.Success = true

		// check if collectionRights already exists
		collectionRightExistingBytes, err := stub.GetState(collectionRight.CollectionRightUUID)
		if collectionRightExistingBytes == nil {
			collectionRightsResponse.Success = false
			collectionRightsResponse.Message = "Collection right does not exist!"
			collectionRightsResponses = append(collectionRightsResponses, collectionRightsResponse)
			collectionRightsOutput.FailureCount++
			continue
		}

		// convert royalty statement to bytes
		collectionRightBytes, err := objectToJSON(collectionRight)
		if err != nil {
			collectionRightsResponse.Success = false
			collectionRightsResponse.Message = err.Error()
			collectionRightsResponses = append(collectionRightsResponses, collectionRightsResponse)
			collectionRightsOutput.FailureCount++
			continue
		}

		// add royalty statement to the ledger
		err = stub.PutState(collectionRight.CollectionRightUUID, collectionRightBytes)
		if err != nil {
			collectionRightsResponse.Success = false
			collectionRightsResponse.Message = err.Error()
		}

		if collectionRightsResponse.Success {
			collectionRightsOutput.SuccessCount++
		} else {
			collectionRightsResponses = append(collectionRightsResponses, collectionRightsResponse)
			collectionRightsOutput.FailureCount++
		}
	}

	collectionRightsOutput.CollectionRightsResponses = collectionRightsResponses

	objBytes, _ := objectToJSON(collectionRightsOutput)
	logger.Info("EXITING <", methodName)
	return shim.Success(objBytes)
}
