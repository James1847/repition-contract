/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-samples/chaincode/repition-contract/chaincode"
)

func main() {
	dataChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating Data chaincode: %v", err)
	}

	if err := dataChaincode.Start(); err != nil {
		log.Panicf("Error starting data chaincode: %v", err)
	}
}