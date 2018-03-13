/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the car structure, with 4 properties.  Structure tags are used by encoding/json library
type Candy struct {
	Name    string `json:"name"`
	Texture string `json:"texture"`
	Colour  string `json:"colour"`
	Owner   string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "fabcar" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabcar"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryCandy" {
		return s.queryCandy(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createCandy" {
		return s.createCandy(APIstub, args)
	} else if function == "queryAllCandies" {
		return s.queryAllCandies(APIstub)
	} else if function == "changeCandyOwner" {
		return s.changeCandyOwner(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryCandy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	carAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(carAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	candies := []Candy{
		Candy{Name: "JollyRancher", Texture: "Hard Candy", Colour: "Multi", Owner: "A"},
		Candy{Name: "Snickers", Texture: "Chewy", Colour: "Brown", Owner: "B"},
		Candy{Name: "KitKat", Texture: "Crunchy", Colour: "Brown", Owner: "A"},
		Candy{Name: "SourPatch", Texture: "Chewy", Colour: "Multi", Owner: "B"},
		Candy{Name: "Kisses", Texture: "Creamy", Colour: "Brown", Owner: "A"},
		Candy{Name: "Hersheys", Texture: "Creamy", Colour: "Brown", Owner: "A"},
		Candy{Name: "DumDums", Texture: "Hard Candy", Colour: "Multi", Owner: "B"},
	}

	i := 0
	for i < len(candies) {
		fmt.Println("i is ", i)
		candyAsBytes, _ := json.Marshal(candies[i])
		APIstub.PutState("CANDY"+strconv.Itoa(i), candyAsBytes)
		fmt.Println("Added", candies[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createCandy(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var candy = Candy{Name: args[1], Texture: args[2], Colour: args[3], Owner: args[4]}

	candyAsBytes, _ := json.Marshal(candy)
	APIstub.PutState(args[0], candyAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllCandies(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "CANDY0"
	endKey := "CANDY999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllCandies:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeCandyOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	candyAsBytes, _ := APIstub.GetState(args[0])
	candy := Candy{}

	json.Unmarshal(candyAsBytes, &candy)
	candy.Owner = args[1]

	candyAsBytes, _ = json.Marshal(candy)
	APIstub.PutState(args[0], candyAsBytes)

	return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
