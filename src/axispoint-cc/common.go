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
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Response -  Object to store Response Status and Message
// ================================================================================
type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

/////////////////////////////////////////////////////
// Constant for table names
/////////////////////////////////////////////////////
const (
	ROYALTYREPORT string = "ROYALTYREPORT"
)

// jsonToObject - common function for unmarshalls : jsonToObject function unmarshalls a JSON into an object
// ================================================================================
func jsonToObject(data []byte, object interface{}) error {
	if err := json.Unmarshal([]byte(data), object); err != nil {
		logger.Errorf("Unmarshal failed : %s ", err.Error())
		return err
	}
	return nil
}

// objectToJSON - common function for marshalls : objectToJSON function marshalls an object into a JSON
// ================================================================================
func objectToJSON(object interface{}) ([]byte, error) {
	var byteArray []byte
	var err error

	if byteArray, err = json.Marshal(object); err != nil {
		logger.Errorf("Marshal failed : %s ", err.Error())
		return nil, err
	}

	if len(byteArray) == 0 {
		return nil, fmt.Errorf(("failed to convert object"))
	}
	return byteArray, nil
}

// getSuccessResponse - Create Success Response and return back to the calling application
// ================================================================================
func getSuccessResponse(message string) pb.Response {
	objResponse := Response{Status: "200", Message: message}
	logger.Info("getSuccessResponse: Called For: ", objResponse)
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
	logger.Info("getErrorResponse: Called For: ", objResponse)
	response, err := json.Marshal(objResponse)
	if err != nil {
		logger.Errorf(fmt.Sprintf("Invalid function %s", err))
	}
	return shim.Success(response)
}

// getHistory - Get History of an asset
// ================================================================================
func getHistory(stub shim.ChaincodeStubInterface, objectType string, args []string) pb.Response {
	var history []AuditHistory

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting one or more keys")
	}

	Key, _ := stub.CreateCompositeKey(objectType, args)
	fmt.Printf("- start getHistory: %s\n", Key)

	// Get History
	resultsIterator, err := stub.GetHistoryForKey(Key)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		historyData, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		var tx AuditHistory
		tx.TxID = historyData.TxId
		tx.Value = string(historyData.Value)
		tx.TimeStamp = historyData.GetTimestamp().String()

		history = append(history, tx) //add this tx to the list
	}
	fmt.Printf("- getHistory returning:\n%s", history)

	//change to array of bytes
	historyAsBytes, _ := json.Marshal(history) //convert to array of bytes
	return shim.Success(historyAsBytes)
}
