package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/*
* addIpiOrg function contains business logic to insert a mapping between IPI and Org to the Ledger
*
* @params   {Array} args
* @property {string} 0       - IPI-Org object
* @return   {pb.Response}    - peer Response
 */
func addIpiOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addIpiOrg"
	logger.Info("ENTERING >", methodName, args)

	// check if array length is greater 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: IPI-Org Mapping object is required")
	}

	ipiOrg := IpiOrgMap{}

	// unmarshal the args input to an IpiOrg struct
	err := jsonToObject([]byte(args[0]), &ipiOrg)
	if err != nil {
		return getErrorResponse(err.Error())
	}

	ipiOrg.DocType = IPIORGMAP
	ipiOrgKey := ipiOrg.Ipi
	//Checking the ledger to confirm that the mapping doesn't exist
	prevIpiOrg, _ := stub.GetState(ipiOrgKey)
	if prevIpiOrg != nil {
		errorMessage := "IPI-Org mapping already exists with this key: " + ipiOrgKey
		logger.Error(methodName, errorMessage)
		return shim.Error(errorMessage)
	}

	byteVal, _ := objectToJSON(ipiOrg)
	err = stub.PutState(ipiOrgKey, byteVal)
	if err != nil {
		errorMessage := "Error committing data for key: " + ipiOrgKey
		logger.Error(methodName, errorMessage)
		return shim.Error(errorMessage)
	}

	logger.Info("EXITING <", methodName, ipiOrg)
	return shim.Success(nil)
}
