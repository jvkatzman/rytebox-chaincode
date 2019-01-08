package main

import (
	"fmt"
	"reflect"
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
		collectionRight.DocType = COLLECTIONRIGHTREPORT
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

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\"}}", COLLECTIONRIGHTREPORT)
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
		collectionRight.DocType = COLLECTIONRIGHTREPORT
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

//generateCollectionStatement -- generate statement for collection or ownership
func generateCollectionStatement(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "generateCollectionStatement"
	logger.Info("ENTERING >", methodName, args)
	var errMessage string

	if len(args) < 3 {
		errMessage = fmt.Sprintf("%s - Incorrect number of parameters provided '%d'.  Operation cannot continue", methodName, len(args))
		logger.Error(errMessage)
		return getErrorResponse(errMessage)
	}
	logger.Infof("%s - parameters received: %s", methodName, strings.Join(args, ","))
	//expUUID, TargeIPI, type
	exploitationReport := ExploitationReport{}
	previousRoyaltyStatement := RoyaltyStatement{}
	royaltyStatementUUID := args[0]
	targetIPI := args[1]
	collectionType := args[2]
	royaltyStatement := RoyaltyStatement{}

	// get the royalty statement report with the UUID exists on the ledger.
	royaltyStatementExistingBytes, err := stub.GetState(royaltyStatementUUID)
	if err != nil {
		errMessage = fmt.Sprintf("%s - Failed to get royalty statement  with uuid '%s' from the ledger.  Error: %s", methodName, royaltyStatementUUID, err.Error())
		logger.Error(errMessage)
		return getErrorResponse(errMessage)
	}
	err = jsonToObject(royaltyStatementExistingBytes, &previousRoyaltyStatement)
	if err != nil {
		errMessage = fmt.Sprintf("%s - Failed to convert royalty statement with uuid '%s'.  Error: %s", methodName, royaltyStatementUUID, err.Error())
		logger.Error(errMessage)
		return getErrorResponse(errMessage)
	}
	expReportUUID := previousRoyaltyStatement.ExploitationReportUUID

	// get the exploitation report with the UUID exists on the ledger.
	exploitationReportExistingBytes, err := stub.GetState(expReportUUID)
	if err != nil {
		errMessage = fmt.Sprintf("%s - Failed to get exploitation report with uuid '%s' from the ledger.  Error: %s", methodName, expReportUUID, err.Error())
		logger.Error(errMessage)
		return getErrorResponse(errMessage)
	}

	err = jsonToObject(exploitationReportExistingBytes, &exploitationReport)
	if err != nil {
		errMessage = fmt.Sprintf("%s - Failed to convert exploitation reqort with uuid '%s'.  Error: %s", methodName, expReportUUID, err.Error())
		logger.Error(errMessage)
		return getErrorResponse(errMessage)
	}

	// create exploitation report parameters to evaluate the selector expressions
	exploitationReportParameters, _ := getEvaluableParameters(&exploitationReport)

	//1. setup the base royalty statement
	royaltyStatement.DocType = ROYALTYSTATEMENT
	royaltyStatement.Source = exploitationReport.Source
	royaltyStatement.SongTitle = exploitationReport.SongTitle
	royaltyStatement.Isrc = exploitationReport.Isrc
	royaltyStatement.ExploitationReportUUID = exploitationReport.ExploitationReportUUID
	royaltyStatement.ExploitationDate = exploitationReport.ExploitationDate
	royaltyStatement.WriterName = exploitationReport.WriterName
	royaltyStatement.Units = exploitationReport.Units
	royaltyStatement.Territory = exploitationReport.Territory
	royaltyStatement.UsageType = exploitationReport.UsageType
	royaltyStatement.Administrator = ""
	royaltyStatement.Collector = ""
	royaltyStatement.Amount = previousRoyaltyStatement.Amount
	logger.Infof("%s - struct value : %+v\n", methodName, royaltyStatement)

	//2. Get all potential collectionrights  that we have based on the 'From' field matching the target IPI.
	collectionRights, err := getCollectionRightsMatchingIpi(stub, targetIPI)
	if err != nil {
		errMessage = fmt.Sprintf("%s - Failed to get collection rights matching target IPI '%s'.  Error: %s", methodName, targetIPI, err.Error())
		logger.Error(errMessage)
		return getErrorResponse(errMessage)
	}

	isSelectorValid := false

	//3.  evaluate every single selector within [collectionRights.rightHolder] until we find the 'FIRST' one that works
	//we should idealy only get one match
collectionRightsLoop:
	for _, collectionRight := range collectionRights {
		//evaluate the rule
		for _, rightHolder := range collectionRight.RightHolders {

			if rightHolder.Selector == "" {
				isSelectorValid = true
			} else {
				isSelectorValidResult, err := evaluate(rightHolder.Selector, exploitationReportParameters)
				if err != nil {
					logger.Errorf("%s - Failed to get a valid evaluator for right holder ipi %s, selector %s. Error: %s", methodName, rightHolder.IPI, rightHolder.Selector, err.Error())
				}
				//to prevent a crash when test cases are run.
				if reflect.ValueOf(isSelectorValidResult).IsValid() {
					isSelectorValid = isSelectorValidResult.(bool)
				}
			}
			if isSelectorValid == true {
				if collectionType == OWNERSHIP {
					royaltyStatement.RightHolder = targetIPI
					royaltyStatement.Administrator = rightHolder.IPI
					royaltyStatement.RightType = COLLECTION
					royaltyStatement.CollectionRightPercent = rightHolder.Percent / 100
					royaltyStatement.CollectionRight = royaltyStatement.Amount * royaltyStatement.CollectionRightPercent
				}
				if collectionType == COLLECTION {
					royaltyStatement.Administrator = targetIPI
					royaltyStatement.Collector = rightHolder.IPI
					royaltyStatement.RightType = COLLECTION
					royaltyStatement.CollectionRightPercent = rightHolder.Percent / 100
					royaltyStatement.CollectionRight = royaltyStatement.Amount * royaltyStatement.CollectionRightPercent
					royaltyStatement.RightHolder = previousRoyaltyStatement.RightHolder
				}

				break collectionRightsLoop
			}
		}
	}
	//we found atleast 1 matching rule above so return the royaltyStatement from above step.
	if isSelectorValid == false {
		//Generate last record
		royaltyStatement.Collector = targetIPI
		royaltyStatement.CollectionRightPercent = 0.0000
		royaltyStatement.CollectionRight = 0.0000
		royaltyStatement.RightType = COLLECTION
		royaltyStatement.RightHolder = previousRoyaltyStatement.RightHolder
	}
	//return the royaltyStatement
	objResultBytes, err := objectToJSON(royaltyStatement)
	if err != nil {
		errMessage = fmt.Sprintf("%s - Failed to get royalty statement bytes for uuid %s.  Error: %s", methodName, royaltyStatement.RoyaltyStatementUUID, err.Error())
		logger.Error(errMessage)
		return getErrorResponse(errMessage)
	}

	return shim.Success(objResultBytes)
}

//getCollectionRightsMatchingIpi Get all potential collectionrights  that we have based on the 'From' field matching the target IPI.
func getCollectionRightsMatchingIpi(stub shim.ChaincodeStubInterface, targetIPI string) ([]CollectionRight, error) {

	var methodName = "getCollectionRightsMatchingIpi"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - target ipi received : %s", methodName, targetIPI)
	defer logger.Infof("%s - End Execution ", methodName)

	queryString := fmt.Sprintf(`{"selector":{"docType":"%s","from":"%s"}}`, COLLECTIONRIGHTREPORT, targetIPI)

	logger.Infof("%s - executing rich query : %s.", methodName, queryString)

	// get royalty statements based on the rich query selector
	queryResult, err := getCollectionRightsForQueryString(stub, queryString)
	if err != nil {
		return nil, fmt.Errorf("%s - Failed to get collection rights with matching target IPI '%s'.  Error: %s", methodName, targetIPI, err.Error())
	}

	var resultCollectionRights []CollectionRight
	err = sliceToStruct(queryResult, &resultCollectionRights)
	if err != nil {
		return nil, fmt.Errorf("%s - Failed to convert bytes to collection rights object for target IPI '%s' .  Error: %s", methodName, targetIPI, err.Error())
	}
	logger.Infof("%s - retrieved %d collection right reports matching target IPI '%s'.", methodName, len(resultCollectionRights), targetIPI)
	return resultCollectionRights, nil

}
