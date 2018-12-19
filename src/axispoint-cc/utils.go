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
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/Knetic/govaluate"

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
	return getSuccessResponse(fmt.Sprintf("resetLedger - deleted %d records.", recordsDeletedCount))
}

// return a default ping response
func ping(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	logger.Infof("Chaincode pinged successfully..")
	return getSuccessResponse("Ping OK")
}

// DeleteAsset - Delete asset based on docType
func deleteAsset(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "deleteAsset"
	logger.Info("ENTERING >", methodName, args)

	recordsDeletedCount := 0
	for _, arg := range args {

		resultIterator, err := stub.GetQueryResult(fmt.Sprintf(`{"selector": {"docType": "%s"}}`, arg))

		if err != nil {
			return getErrorResponse(err.Error())
		}

		defer resultIterator.Close()

		for resultIterator.HasNext() {
			result, err := resultIterator.Next()

			err = stub.DelState(result.Key)
			if err != nil {
				return getErrorResponse(err.Error())
			}
			recordsDeletedCount++
		}
	}

	logger.Info("EXITING <", methodName)
	return getSuccessResponse(fmt.Sprintf("deleteAsset - deleted %d records.", recordsDeletedCount))
}

// DeleteAsset - Delete asset based on UUIDs
func deleteAssetByUUID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "deleteAssetByUUID"
	logger.Info("ENTERING >", methodName, args)

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: UUID is missing")
	}

	recordsDeletedCount := 0
	for _, arg := range args {

		err := stub.DelState(arg)
		if err != nil {
			return getErrorResponse(err.Error())
		}
		recordsDeletedCount++

	}

	logger.Info("EXITING <", methodName)
	return getSuccessResponse(fmt.Sprintf("deleteAssetByUUID - deleted %d records.", recordsDeletedCount))
}

// getAssetByUUID - Get asset based on UUID
func getAssetByUUID(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var methodName = "getAssetByUUID"
	logger.Info("ENTERING >", methodName, args)

	//Check if array length is greater than 0
	if len(args) < 1 {
		return getErrorResponse("Missing arguments: UUID is missing")
	}

	objectBytes, err := stub.GetState(args[0])
	if err != nil {
		return getErrorResponse(err.Error())
	}
	if objectBytes == nil {
		return getErrorResponse(fmt.Sprintf("UUID: %s does not exist", args[0]))
	}

	//return bytes as result
	return shim.Success(objectBytes)
}

//getEvaluableParameters - Returns the struct parameters using reflect
func getEvaluableParameters(asset interface{}) (map[string]interface{}, error) {
	if asset == nil || reflect.ValueOf(asset).Kind() != reflect.Ptr || reflect.ValueOf(asset).Elem().Kind() != reflect.Struct {
		err := errors.New("wrong parameter type: the 'asset' param has to be a struct instance pointer")
		logger.Error(err)
		return nil, err
	}

	val := reflect.ValueOf(asset).Elem()
	logger.Infof("getEvaluableParameters or struct %s", val.Type())

	parameters := make(map[string]interface{}, 8)
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		parameters[typeField.Name] = valueField.Interface()
	}

	return parameters, nil
}

//evaluate - Evaluates the selector againist the parameters and returns true/false
func evaluate(selector string, parameters map[string]interface{}) (interface{}, error) {
	logger.Infof("evaluate selector: %s", selector)
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

// inTimeRange - Return true/false if the target datetime range is with in the specifed datetime range
func inTimeRange(day1, day2, day3, day4 string) bool {
	t1, _ := time.Parse(time.RFC3339, day1) //start
	t2, _ := time.Parse(time.RFC3339, day2) //end
	t3, _ := time.Parse(time.RFC3339, day3) //target.start
	t4, _ := time.Parse(time.RFC3339, day4) //target.end

	return (((t1.Before(t3) || t1.Equal(t3)) && (t2.After(t3) || t2.Equal(t3))) ||
		((t1.Before(t4) || t1.Equal(t4)) && (t2.After(t4) || t2.Equal(t4))) ||
		((t1.After(t3) || t1.Equal(t3)) && (t2.Before(t4) || t2.Equal(t4))))
}

// isDateInTimeRange - Return true/false if the target datetime is with in the specifed datetime range
func isDateInTimeRange(day1, day2, day3 string) bool {
	t1, _ := time.Parse(time.RFC3339, day1) //target
	t2, _ := time.Parse(time.RFC3339, day2) //start
	t3, _ := time.Parse(time.RFC3339, day3) //end

	return (t1.After(t2) || t1.Equal(t2)) && (t1.Before(t3) || t1.Equal(t3))
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
