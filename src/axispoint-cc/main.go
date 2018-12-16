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
	"fmt"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// AxispointChaincode implementation
type AxispointChaincode struct {
	tableMap map[string]int
	funcMap  map[string]InvokeFunc
}

var logger = shim.NewLogger("axispoint-cc")

type InvokeFunc func(stub shim.ChaincodeStubInterface, args []string) pb.Response

// initFunctionMaps - Map all the Functions here for Invoke
/////////////////////////////////////////////////////
func (t *AxispointChaincode) initFunctionMaps() {
	t.tableMap = make(map[string]int)
	t.funcMap = make(map[string]InvokeFunc)
	t.funcMap["addRoyaltyStatements"] = addRoyaltyStatements
	t.funcMap["addExploitationReports"] = addExploitationReports
	t.funcMap["updateExploitationReports"] = updateExploitationReports
	t.funcMap["getExploitationReports"] = getExploitationReports
	t.funcMap["getRoyaltyStatements"] = getRoyaltyStatements
	t.funcMap["resetLedger"] = resetLedger
	t.funcMap["ping"] = ping
	t.funcMap["addCopyrightDataReports"] = addCopyrightDataReports
	t.funcMap["getCopyrightDataReportByID"] = getCopyrightDataReportByID
	t.funcMap["deleteCopyrightDataReportByIDs"] = deleteCopyrightDataReportByIDs
	t.funcMap["updateCopyrightDataReports"] = updateCopyrightDataReports
	t.funcMap["searchForCopyrightDataReportWithParameters"] = searchForCopyrightDataReportWithParameters
	t.funcMap["getAllCopyrightDataReports"] = getAllCopyrightDataReports
	t.funcMap["deleteAsset"] = deleteAsset
	t.funcMap["deleteAssetByUUID"] = deleteAssetByUUID
	t.funcMap["getAssetByUUID"] = getAssetByUUID
	t.funcMap["addOwnerAdministrations"] = addOwnerAdministrations
	t.funcMap["updateOwnerAdministrations"] = updateOwnerAdministrations
	t.funcMap["getOwnerAdministrations"] = getOwnerAdministrations
	t.funcMap["addAdministratorAffiliations"] = addAdministratorAffiliations
	t.funcMap["updateAdministratorAffiliations"] = updateAdministratorAffiliations
	t.funcMap["getAdministratorAffiliations"] = getAdministratorAffiliations
	t.funcMap["getRoyaltyStatementsByUUIDs"] = getRoyaltyStatementsByUUIDs
	t.funcMap["updateRoyaltyStatements"] = updateRoyaltyStatements
	t.funcMap["addReports"] = addReports
}
func addReports(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var methodName = "addReports"
	logger.Infof("%s - Begin Execution ", methodName)
	logger.Infof("%s - # of arguments: %d", methodName, len(args))
	logger.Infof("%s - parameters received : %s", methodName, strings.Join(args, ","))
	defer logger.Infof("%s - End Execution ", methodName)

	if len(args) < 2 {
		return getErrorResponse(methodName + " incorrect number of args provided.")
	}

	//args = append(args[:0], args[1:]...)
	result := insertExploitationReports(stub, args)
	if result.GetStatus() != shim.OK {
		return getErrorResponse("updateExploitationReports failed with " + result.GetMessage())
	}
	args = append(args[:0], args[1:]...)
	result = addRoyaltyStatements(stub, args)
	if result.GetStatus() != shim.OK {
		return getErrorResponse("addRoyaltyStatements failed with " + result.GetMessage())
	}
	return shim.Success(nil)
}

// Init - intialize chaincode
func (t *AxispointChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	// Initialize the chaincode
	logger.Info("########### AxispointChaincode Init ###########")
	t.initFunctionMaps()
	isInit = true

	return shim.Success(nil)
}

// Invoke - invoke/query on chaincode
func (t *AxispointChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	if !isInit {
		t.initFunctionMaps()
		isInit = true
	}
	logger.Info("########### Invoke/Query ###########")
	function, args := stub.GetFunctionAndParameters()

	f, ok := t.funcMap[function]
	if ok {
		return f(stub, args)
	}

	logger.Errorf("Invalid function name %s", function)
	return getErrorResponse(fmt.Sprintf("Invalid function %s", function))
}

var isInit = false

func main() {
	err := shim.Start(new(AxispointChaincode))
	if err != nil {
		fmt.Printf("Error starting Axispoint chaincode: %s", err)
	}
}
