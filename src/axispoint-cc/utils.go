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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Knetic/govaluate"
	"reflect"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// getObjectByQueryFromLedger - Get all objects matching object type from the ledger.
// ============================================================
func getObjectByQueryFromLedger(stub shim.ChaincodeStubInterface, query string) ([]string, error) {
	var methodName = "getObjectByQueryFromLedger"
	logger.Info("ENTERING >", methodName, query)

	resultIterator, err := stub.GetQueryResult(query)
	if err != nil {
		errorMessage := "GetQueryResult - error: " + err.Error()
		logger.Error(methodName, errorMessage)
		return nil, errors.New(errorMessage)
	}

	defer resultIterator.Close()

	slice := make([]string, 0) // be nice and return at least an empty slice
	for resultIterator.HasNext() {
		result, err := resultIterator.Next()
		if err != nil {
			errorMessage := "Iterator failed - error: " + err.Error()
			logger.Error(methodName, errorMessage)
			return nil, errors.New(errorMessage)
		}

		slice = append(slice, string(result.Value))
	}

	logger.Info("EXITING <", methodName, slice)
	return slice, nil
}

// jsonToObject - common function for unmarshalls : jsonToObject function unmarshalls a JSON into an object
// ================================================================================
func jsonToObject(data []byte, object interface{}) error {
	var methodName = "jsonToObject"
	logger.Info("ENTERING >", methodName)
	if err := json.Unmarshal([]byte(data), object); err != nil {
		errorMessage := "Unmarshal failed - error: " + err.Error()
		logger.Error(methodName, errorMessage)
		return errors.New(errorMessage)
	}
	logger.Info("EXITING <", methodName, object)
	return nil
}

// objectToJSON - common function for marshalls : objectToJSON function marshalls an object into a JSON
// ================================================================================
func objectToJSON(object interface{}) ([]byte, error) {
	var methodName = "objectToJSON"
	logger.Info("ENTERING >", methodName, object)
	var byteArray []byte
	var err error

	if byteArray, err = json.Marshal(object); err != nil {
		errorMessage := "Marshal failed - error: " + err.Error()
		logger.Error(methodName, errorMessage)
		return nil, errors.New(errorMessage)
	}

	if len(byteArray) == 0 {
		errorMessage := "failed to convert object"
		logger.Error(methodName, errorMessage)
		return nil, errors.New(errorMessage)
	}
	logger.Info("EXITING <", methodName)
	return byteArray, nil
}

// sliceToStruct - unmarshals a []string into a the given []object instance.
// ============================================================
func sliceToStruct(slice []string, object interface{}) error {
	var methodName = "sliceToStruct"
	logger.Info("ENTERING >", methodName, slice, object)

	jsonArray := sliceToJSONArray(slice)
	err := json.Unmarshal([]byte(jsonArray), &object)
	if err != nil {
		errorMessage := "Unmarshal failed - error: " + err.Error()
		logger.Error(methodName, errorMessage)
		return errors.New(errorMessage)
	}

	logger.Info("EXITING <", methodName, object)
	return nil
}

// sliceToJSONArray - Produces a JSON array out of a slice of strings
// ============================================================
func sliceToJSONArray(slice []string) string {
	var methodName = "sliceToJSONArray"
	logger.Info("ENTERING >", methodName, slice)

	jsonString := fmt.Sprintf("%s", slice)
	jsonString = strings.Replace(jsonString, `} {`, `}, {`, -1)

	logger.Info("EXITING <", methodName, jsonString)
	return jsonString
}

// getSuccessResponse - Create Success Response and return back to the calling application
// ================================================================================
func getSuccessResponse(message string) pb.Response {
	objResponse := Response{Status: "200", Message: message}
	response, err := json.Marshal(objResponse)
	if err != nil {
		logger.Errorf(fmt.Sprintf("Invalid function %s", err))
	}
	return shim.Success(response)
}

// getErrorResponse - Create Error Response and return back to the calling application
// ================================================================================
func getErrorResponse(message string) pb.Response {
	objResponse := Response{Status: "500", Message: message}
	response, err := json.Marshal(objResponse)
	if err != nil {
		logger.Errorf(fmt.Sprintf("Invalid function %s", err))
	}
	return shim.Success(response)
}

// resetWorldState - remove all data from the world state
// ================================================================================
func resetWorldState(stub shim.ChaincodeStubInterface) (int, error) {
	methodName := "resetWorldState"
	logger.Infof("Begin execution - %s.", methodName)
	defer logger.Infof("End execution - %s.", methodName)
	startKey := ""
	endKey := ""
	recordsDeletedCount := 0

	iterator, err := stub.GetStateByRange(startKey, endKey)
	defer iterator.Close()
	if err != nil {
		message := fmt.Sprintf("%s - Failed to get state by range with error: %s", methodName, err)
		logger.Error(message)
		return recordsDeletedCount, err
	}

	for iterator.HasNext() {
		responseRange, err := iterator.Next()
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to get next record from iterator: %s", err.Error())
			logger.Error(errorMsg)
			return recordsDeletedCount, err
		}

		recordKey := responseRange.GetKey()
		logger.Infof("About to delete record with key %s", recordKey)
		err = stub.DelState(recordKey)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to delete record '%d' with key %s: %s", recordsDeletedCount, recordKey, err.Error())
			logger.Error(errorMsg)
			return recordsDeletedCount, err
		}
		recordsDeletedCount++
		logger.Debugf("%s - Successfully deleted record '%d' with key: %s", methodName, recordsDeletedCount, recordKey)
	}
	logger.Infof("%s - Total # of records deleted : %d", methodName, recordsDeletedCount)
	return recordsDeletedCount, nil
}

//resetLedger - remove all data from the world state.
/*
* @params   {Array} args - empty array
* @return   {pb.Response}    - peer Response
 */
func resetLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	recordsDeletedCount, err := resetWorldState(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(fmt.Sprintf("resetLedger - deleted %d records.", recordsDeletedCount)))
}

// return a default ping response
func ping(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Infof("Chaincode pinged successfully..")
	return shim.Success([]byte("Ping OK"))
}

/**
 * Evaluates a parametrized boolean expression. The expression's parameters are sent in a struct instance.
 *
 * @param {string} selector - the boolean parametrize expression to evaluate
 * @param {&struct} asset - a pointer to any struct instance
 * @return {(boolean, error)} (boolean, error) - the
 */
func evaluate(selector string, asset interface{}) (interface{}, error) {
	if asset == nil {
		err := errors.New("wrong parameter value: the 'asset' param cannot be nil")
		logger.Error(err)
		return nil, err
	}

	if reflect.ValueOf(asset).Kind() != reflect.Ptr {
		err := errors.New("wrong parameter type: the 'asset' param has to be a struct pointer")
		logger.Error(err)
		return nil, err
	}

	if reflect.ValueOf(asset).Elem().Kind() != reflect.Struct {
		err := errors.New("wrong parameter type: the 'asset' param has to be a struct pointer")
		logger.Error(err)
		return nil, err
	}

	val := reflect.ValueOf(asset).Elem()
	logger.Infof("evaluate selector: %s for struct %s", selector, val.Type())

	parameters := make(map[string]interface{}, 8)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		parameters[typeField.Name] = valueField.Interface()
	}

	expression, err := govaluate.NewEvaluableExpression(selector)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	result, err := expression.Evaluate(parameters)
	if err != nil {
		logger.Error(err)
	}

	return result, err
}
