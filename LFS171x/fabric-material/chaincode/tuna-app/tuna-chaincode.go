// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
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

/* Define Tuna structure, with 4 properties.  
Structure tags are used by encoding/json library
*/
type Tuna struct {
	Clarity string `json:"clarity"`
	Color string `json:"color"`
	Cut  string `json:"cut"`
	Carat  string `json:"carat"`
	Certification  string `json:"cert"`
	Name string `json:"name"`
	TransId string `json:"transid"`
	Holder string `json:"holdername"`
	TimeStamp string `json:"timeStamp"`
	Type string `json:"type"`
	Image string `json:"image"`
	Latitude string `json:"latitude"`
	Longitude string `json:"longitude"`
}

/*
 * The Init method *
 called when the Smart Contract "tuna-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "tuna-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryTunaHistory" {
		return s.queryTunaHistory(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordTuna" {
		return s.recordTuna(APIstub, args)
	} else if function == "queryAllTuna" {
		return s.queryAllTuna(APIstub)
	} else if function == "changeTunaHolder" {
		return s.changeTunaHolder(APIstub, args)
	}else if function == "updateLatLong"{
		return s.updateLatLong(APIstub,args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryTuna method *
Used to view the records of one particular tuna
It takes one argument -- the key for the tuna in question
 */
func (s *SmartContract) queryTuna(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// if len(args) != 1 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 1")
	// }

	// tunaAsBytes, _ := APIstub.GetState(args[0])
	// if tunaAsBytes == nil {
	// 	return shim.Error("Could not locate tuna")
	// }
	// return shim.Success(tunaAsBytes)


	fmt.Println("Entering Query Food information")

    // Assuming food key is at zero index
    historyIer, err := APIstub.GetHistoryForKey(args[0])

    if err != nil {
        errMsg := fmt.Sprintf("[ERROR] cannot retrieve history of food record with id <%s>, due to %s", args[0], err)
        fmt.Println(errMsg)
        return shim.Error(errMsg)
    }

    result := make([]Tuna, 0)
    for historyIer.HasNext() {
        modification, err := historyIer.Next()
        if err != nil {
            errMsg := fmt.Sprintf("[ERROR] cannot read food record modification, id <%s>, due to %s", args[0], err)
            fmt.Println(errMsg)
            return shim.Error(errMsg)
        }
        var food Tuna
        json.Unmarshal(modification.Value, &food)
        result = append(result, food)
    }

    outputAsBytes, _ := json.Marshal(&result)                   
    return shim.Success(outputAsBytes)
}



/*
 * The initLedger method *
Will add test data (10 tuna catches)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	tuna := []Tuna{
		Tuna{Clarity: "913F", Color: "red", Cut: "good", Carat: "1",Certification:"IGI",Name: "Red Diamond",TransId:"dvbadjhbvdvadb7bvwebvuwvwbvuwebbvewbuew84be",Holder:"Akshay",TimeStamp:"", Type: "Add",Image:"",Latitude:"18.45",Longitude:"73.565"},
		Tuna{Clarity: "913F", Color: "blue", Cut: "good", Carat: "1",Certification:"HRD",Name: "Blue Diamond",TransId:"dgkhwr7h4iubg37g3b4ubge83b4u7ewurjdqw6te26",Holder:"Akash",TimeStamp:"", Type: "Add",Image:"",Latitude:"18.45",Longitude:"73.565"}, 
		Tuna{Clarity: "913F", Color: "yellow", Cut: "good", Carat: "1",Certification:"GIA",Name: "Yellow Diamond",TransId:"evfb2734ghi3hgubg28hg82h4gg8nhg82g47fy432f",Holder:"Kapil",TimeStamp:"", Type: "Add",Image:"",Latitude:"18.45",Longitude:"73.565"}, 

	}

	i := 0
	for i < len(tuna) {
		fmt.Println("i is ", i)
		tunaAsBytes, _ := json.Marshal(tuna[i])
		APIstub.PutState(strconv.Itoa(i+1), tunaAsBytes)
		fmt.Println("Added", tuna[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordTuna method *
Fisherman like Sarah would use to record each of her tuna catches. 
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordTuna(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 14 {
		return shim.Error("Incorrect number of arguments. Expecting 14")
	}

	var tuna = Tuna{Clarity: args[1], Color: args[2], Cut: args[3], Carat: args[4] , Certification: args[5], Name: args[6],TransId: args[7],Holder: args[8],TimeStamp:args[9], Type:args[10],Image: args[11],Latitude: args[12],Longitude: args[13]}

	tunaAsBytes, _ := json.Marshal(tuna)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record tuna catch: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllTuna method *
allows for assessing all the records added to the ledger(all tuna catches)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllTuna(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999"

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
		// Add comma before array members,suppress it for the first array member
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

	fmt.Printf("- queryAllTuna:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changeTunaHolder method *
The data in the world state can be updated with who has possession. 
This function takes in 2 arguments, tuna id and new holder name. 
 */
func (s *SmartContract) changeTunaHolder(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	tunaAsBytes, _ := APIstub.GetState(args[0])
	if tunaAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	tuna := Tuna{}

	json.Unmarshal(tunaAsBytes, &tuna)
	// Normally check that the specified argument is a valid holder of tuna
	// we are skipping this check for this example
	tuna.Holder = args[1]
	tuna.TransId = args[2]
	tuna.TimeStamp = args[3]
	tuna.Type = args[4]
	tuna.Latitude = args[5]
	tuna.Longitude = args[6]

	tunaAsBytes, _ = json.Marshal(tuna)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change tuna holder: %s", args[0]))
	}

	return shim.Success(nil)
}

func (s *SmartContract) queryTunaHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
    fmt.Println("Entering Query Food information")
    fmt.Sprintf("Hi There")
    // Assuming food key is at zero index
    historyIer, err := APIstub.GetHistoryForKey(args[0])

    if err != nil {
        errMsg := fmt.Sprintf("[ERROR] cannot retrieve history of food record with id <%s>, due to %s", args[0], err)
        fmt.Println(errMsg)
        return shim.Error(errMsg)
    }

    result := make([]Tuna, 0)
    for historyIer.HasNext() {
        modification, err := historyIer.Next()
        if err != nil {
            errMsg := fmt.Sprintf("[ERROR] cannot read food record modification, id <%s>, due to %s", args[0], err)
            fmt.Println(errMsg)
            return shim.Error(errMsg)
        }
        var food Tuna
        json.Unmarshal(modification.Value, &food)
        result = append(result, food)
    }

    outputAsBytes, _ := json.Marshal(&result)                   
    return shim.Success(outputAsBytes)
 }

func (s *SmartContract) updateLatLong(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) !=  5{
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	tunaAsBytes, _ := APIstub.GetState(args[0])
	if tunaAsBytes == nil {
		return shim.Error("Could not locate tuna")
	}
	tuna := Tuna{}

	json.Unmarshal(tunaAsBytes, &tuna)
	// Normally check that the specified argument is a valid holder of tuna
	// we are skipping this check for this example
	tuna.TransId = args[1]
	tuna.TimeStamp = args[2]
	tuna.Latitude = args[3]
	tuna.Longitude = args[4]

	tunaAsBytes, _ = json.Marshal(tuna)
	err := APIstub.PutState(args[0], tunaAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change tuna location: %s", args[0]))
	}

	return shim.Success(nil)
}
 
/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}