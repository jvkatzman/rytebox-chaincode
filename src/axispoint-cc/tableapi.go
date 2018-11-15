/*
* Copyright 2018 IT People Corporation. All Rights Reserved.
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an 'AS IS' BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

////////////////////////////////////////////////////////////////////////////
// Update the Object - Replace current data with replacement
// Register users into this table
////////////////////////////////////////////////////////////////////////////
func UpdateObject(stub shim.ChaincodeStubInterface, objectType string, keys []string, objectData []byte) error {
	// Check number of keys
	err := VerifyAtLeastOneKeyIsPresent(keys)
	if err != nil {
		return err
	}

	// Convert keys to  compound key
	compositeKey, _ := stub.CreateCompositeKey(objectType, keys)

	// Add Object JSON to state
	err = stub.PutState(compositeKey, objectData)
	if err != nil {
		fmt.Printf("UpdateObject() : Error inserting Object into State Database %s", err)
		return err
	}

	return nil

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Retrieve the object based on the key and simply delete it
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DeleteObject(stub shim.ChaincodeStubInterface, objectType string, keys []string) error {
	// Check number of keys
	err := VerifyAtLeastOneKeyIsPresent(keys)
	if err != nil {
		return err
	}

	// Convert keys to  compound key
	compositeKey, _ := stub.CreateCompositeKey(objectType, keys)

	// Remove object from the State Database
	err = stub.DelState(compositeKey)
	if err != nil {
		fmt.Printf("DeleteObject() : Error deleting Object into State Database %s", err)
		return err
	}
	fmt.Println("DeleteObject() : ", "Object : ", objectType, " Key : ", compositeKey)

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Delete all objects of ObjectType
//
////////////////////////////////////////////////////////////////////////////////////////////////////////////
func DeleteAllObjects(stub shim.ChaincodeStubInterface, objectType string) error {
	// Convert keys to  compound key
	compositeKey, _ := stub.CreateCompositeKey(objectType, []string{""})

	// Remove object from the State Database
	err := stub.DelState(compositeKey)
	if err != nil {
		fmt.Printf("DeleteAllObjects() : Error deleting all Object into State Database %s", err)
		return err
	}
	fmt.Println("DeleteAllObjects() : ", "Object : ", objectType, " Key : ", compositeKey)

	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Replaces the Entry in the Ledger
// The existing object is simply queried and the data contents is replaced with
// new content
////////////////////////////////////////////////////////////////////////////////////////////////////////////
func ReplaceObject(stub shim.ChaincodeStubInterface, objectType string, keys []string, objectData []byte) error {
	// Check number of keys
	err := VerifyAtLeastOneKeyIsPresent(keys)
	if err != nil {
		return err
	}

	// Convert keys to  compound key
	compositeKey, _ := stub.CreateCompositeKey(objectType, keys)

	// Add Party JSON to state
	err = stub.PutState(compositeKey, objectData)
	if err != nil {
		fmt.Printf("ReplaceObject() : Error replacing Object in State Database %s", err)
		return err
	}

	fmt.Println("ReplaceObject() : - end init object ", objectType)
	return nil
}

////////////////////////////////////////////////////////////////////////////
// Query a User Object by Object Name and Key
// This has to be a full key and should return only one unique object
////////////////////////////////////////////////////////////////////////////
func QueryObject(stub shim.ChaincodeStubInterface, objectType string, keys []string) ([]byte, error) {
	// Check number of keys
	err := VerifyAtLeastOneKeyIsPresent(keys)
	if err != nil {
		return nil, err
	}

	compoundKey, _ := stub.CreateCompositeKey(objectType, keys)
	fmt.Println("QueryObject() : Compound Key : ", compoundKey)

	objBytes, err := stub.GetState(compoundKey)
	if err != nil {
		return nil, err
	}

	return objBytes, nil
}

////////////////////////////////////////////////////////////////////////////
// Query a User Object by Object Name and Key
// This has to be a full key and should return only one unique object
////////////////////////////////////////////////////////////////////////////
func QueryObjectWithProcessingFunction(stub shim.ChaincodeStubInterface, objectType string, keys []string, fname func(shim.ChaincodeStubInterface, []byte, []string) error) ([]byte, error) {
	// Check number of keys
	err := VerifyAtLeastOneKeyIsPresent(keys)
	if err != nil {
		return nil, err
	}

	compoundKey, _ := stub.CreateCompositeKey(objectType, keys)
	fmt.Println("QueryObject: Compound Key : ", compoundKey)

	objBytes, err := stub.GetState(compoundKey)
	if err != nil {
		return nil, err
	}

	if objBytes == nil {
		return nil, fmt.Errorf("QueryObject: No Data Found for Compound Key : %s", compoundKey)
	}

	// Perform Any additional processing of data
	fmt.Println("fname() : Successful - Proceeding to fname")

	err = fname(stub, objBytes, keys)
	if err != nil {
		jsonResp := "{\"fname() Error\":\" Cannot create Object for key " + compoundKey + "\"}"
		return objBytes, errors.New(jsonResp)
	}

	return objBytes, nil
}

////////////////////////////////////////////////////////////////////////////
// Get a List of Rows based on query criteria from the OBC
// The getList Function
////////////////////////////////////////////////////////////////////////////
func GetKeyList(stub shim.ChaincodeStubInterface, args []string) (shim.StateQueryIteratorInterface, error) {
	// Define partial key to query within objects namespace (objectType)
	objectType := args[0]

	// Check number of keys

	err := VerifyAtLeastOneKeyIsPresent(args[1:])
	if err != nil {
		return nil, err
	}

	// Execute the Query
	// This will execute a key range query on all keys starting with the compound key
	resultsIterator, err := stub.GetStateByPartialCompositeKey(objectType, args[1:])
	if err != nil {
		return nil, err
	}

	defer resultsIterator.Close()

	// Iterate through result set
	var i int
	for i = 0; resultsIterator.HasNext(); i++ {

		// Retrieve the Key and Object
		myCompositeKey, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		fmt.Println("GetList() : my Value : ", myCompositeKey)
	}
	return resultsIterator, nil
}

///////////////////////////////////////////////////////////////////////////////////////////
// GetQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
///////////////////////////////////////////////////////////////////////////////////////////
func GetQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {
	fmt.Printf("GetQueryResultForQueryString() : getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		//buffer.WriteString("{\"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		//buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("GetQueryResultForQueryString(): getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func GetList(stub shim.ChaincodeStubInterface, objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	// Check number of keys
	err := VerifyAtLeastOneKeyIsPresent(keys)
	if err != nil {
		return nil, err
	}

	// Get Result set
	resultIter, err := stub.GetStateByPartialCompositeKey(objectType, keys)
	fmt.Println("GetList(): Retrieving Objects into an array")
	if err != nil {
		return nil, err
	}

	// Return iterator for result set
	// Use code above to retrieve objects
	return resultIter, nil
}

////////////////////////////////////////////////////////////////////////////
// This function verifies if the number of key provided is at least 1 and
// < the max keys defined for the Object
////////////////////////////////////////////////////////////////////////////

func VerifyAtLeastOneKeyIsPresent(args []string) error {
	// Check number of keys
	nKeys := len(args)
	if nKeys == 1 {
		return nil
	}

	if nKeys < 1 {
		err := fmt.Sprintf("VerifyAtLeastOneKeyIsPresent() Failed: Atleast 1 Key must is needed :  nKeys : %s", strconv.Itoa(nKeys))
		fmt.Println(err)
		return errors.New(err)
	}

	return nil
}
