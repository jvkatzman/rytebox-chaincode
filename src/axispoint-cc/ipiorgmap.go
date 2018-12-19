package main

import (
	"errors"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/*
* addIpiOrg function inserts a mapping between IPI and Org to the Ledger
*
* @params   {Array} args
* @property {string} 0       - IPI-Org object
* @return   {pb.Response}    - peer Response
 */
func addIpiOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "addIpiOrg"
	logger.Info("ENTERING >", methodName, args)

	// check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: IPI-Org Mapping object is required")
	}

	err := addUpdateIpiOrg(stub, args[0], false)
	if err != nil {
		errorMessage := err.Error()
		logger.Error(methodName, errorMessage)
		return getErrorResponse(errorMessage)
	}

	logger.Info("EXITING <", methodName)
	return shim.Success(nil)
}

/*
* updateIpiOrg function inserts a new or update an existing mapping between IPI and Org to the Ledger
*
* @params   {Array} args
* @property {string} 0       - IPI-Org object
* @return   {pb.Response}    - peer Response
 */
func updateIpiOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "updateIpiOrg"
	logger.Info("ENTERING >", methodName, args)

	// check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: IPI-Org Mapping object is required")
	}

	err := addUpdateIpiOrg(stub, args[0], true)
	if err != nil {
		errorMessage := err.Error()
		logger.Error(methodName, errorMessage)
		return getErrorResponse(errorMessage)
	}

	logger.Info("EXITING <", methodName)
	return shim.Success(nil)
}

/*
* addUpdateIpiOrg function contains the actual business logic to insert a new or update an existing mapping
* between IPI and Org to the Ledger
*
* @param      {string}       - IPI-Org json object
* @param        {bool}       - updateFlag
* @return      {error}       - Error
 */
func addUpdateIpiOrg(stub shim.ChaincodeStubInterface, ipiOrgObj string, updateFlag bool) error {
	var methodName = "addUpdateIpiOrg"
	logger.Info("ENTERING >", methodName, ipiOrgObj, updateFlag)

	ipiOrg := IpiOrgMap{}

	// unmarshal the args input to an IpiOrg struct
	err := jsonToObject([]byte(ipiOrgObj), &ipiOrg)
	if err != nil {
		return err
	}

	ipiOrg.DocType = IPIORGMAP
	ipiOrgKey := ipiOrg.Ipi

	if !updateFlag {
		//updateFlag==false; This is invoked by a POST request
		//Checking the ledger to confirm that the mapping doesn't exist
		prevIpiOrg, _ := stub.GetState(ipiOrgKey)
		if prevIpiOrg != nil {
			errorMessage := "IPI-Org mapping already exists with this key: " + ipiOrgKey
			logger.Error(methodName, errorMessage)
			return errors.New(errorMessage)
		}
	}

	byteVal, _ := objectToJSON(ipiOrg)
	err = stub.PutState(ipiOrgKey, byteVal)
	if err != nil {
		errorMessage := "Error committing data for key: " + ipiOrgKey
		logger.Error(methodName, errorMessage)
		return errors.New(errorMessage)
	}

	logger.Info("EXITING <", methodName, ipiOrg)
	return nil
}

//getIpiOrgByUUID function retrieves IPI-Org Mappings by IPI (UUID of a participant)
/*
* @params   {Array}  args
* @property {string} 0     - IPI (UUID of a participant)
* @return   {Peer.Reponse} - IPI-Org mapping object as Bytes
 */
func getIpiOrgByUUID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getIpiOrgByUUID"
	logger.Info("ENTERING >", methodName, args)
	return getAssetByUUID(stub, args)

}
