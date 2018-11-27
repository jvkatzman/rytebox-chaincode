package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

//resetLedger - remove all data from the world state.
/*
* @params   {Array} args - empty array
* @return   {pb.Response}    - peer Response
 */
func resetLedger(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	_, err := resetWorldState(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
