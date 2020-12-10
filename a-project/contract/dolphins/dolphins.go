/*
 * Copyright IBM Corp All Rights Reserved
 *
 * SPDX-License-Identifier: Apache-2.0
 */

 package main

 import (
	 "bytes"
	 "encoding/json"
	 "fmt"
	 "strconv"
	 "time"
 
	 "github.com/hyperledger/fabric/core/chaincode/shim"
	 "github.com/hyperledger/fabric/protos/peer"
 )
 
 // SimpleAsset implements a simple chaincode to manage an asset
 type Dolphins struct {
 }
 type Diver struct {
	 Id   string `json:"id"`
	 Name string `json:"name"`
	 Bdate string `json:"bdate"`
	 Gender string `json:"gender"`
	 Btype string `json:"btype"`
	 Levels []Level `json:"levels"`
 }

 type Level struct {
	 Levelname string `json:"levelname"`
	 Org string `json:"org"`
	 Instid string `json:"instid"`
	 Courses []string `json:"courses"`
	 Status string `json:"status"` // incourse -> qualified
 }
 
 // Init is called during chaincode instantiation to initialize any
 // data. Note that chaincode upgrade also calls this function to reset
 // or to migrate data.
 func (t *Dolphins) Init(stub shim.ChaincodeStubInterface) peer.Response {
 
	 return shim.Success(nil)
 }
 
 // Invoke is called per transaction on the chaincode. Each transaction is
 // either a 'get' or a 'add' on the level created by Init function. The add
 // method may create a new diver by specifying a new key-value pair.
 func (t *Dolphins) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	 // Extract the function and args from the transaction proposal
	 fn, args := stub.GetFunctionAndParameters()
 
	 var result string
	 var err error
	 if fn == "addDiver" {
		 result, err = addDiver(stub, args)
	 } else if fn == "addLevel" {
		 result, err = addLevel(stub, args)
	 } else if fn == "addCourse" {
		 result, err = addCourse(stub, args)
	 } else if fn == "addTestResult" {
		result, err = addTestResult(stub, args)
	 } else if fn == "getLevel" {
		result, err = getLevel(stub, args)
	 }  else if fn == "getHistoryForKey" {
		result, err = getHistoryForKey(stub, args)
	 } else {
		 return shim.Error("Not supported chaincode function.")
	 }
 
	 if err != nil {
		 return shim.Error(err.Error())
	 }
 
	 // Return the result as success payload
	 return shim.Success([]byte(result))
 }
 
 // Add diver datas (both key and value) on the ledger. If the key exists,
 // it will override the value with the new one
 func addDiver(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	 if len(args) != 5 {
		 return "", fmt.Errorf("Incorrect arguments. Expecting 5 parameters")
	 }
	 fmt.Println("AddDiver: start")
	 // JSON  변환
	 var data = Diver{Id: args[0], Name: args[1], Bdate: args[2], Gender: args[3], Btype: args[4]}
	 dataAsBytes, _ := json.Marshal(data)
 
	 err := stub.PutState(args[0], dataAsBytes)
	 if err != nil {
		 return "", fmt.Errorf("Failed to add diver: %s", args[0])
	 }
	 return string(dataAsBytes), nil
 }

 func addLevel(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 4 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	// 해당 아이디가 이전에 있는지 검사
    // GetState => return 값이 nil 아니면 중복 아이디
    value, err := stub.GetState(args[0])
    if err != nil {
      return "", fmt.Errorf("Failed to get diver: %s with error: %s", args[0], err)
     }
    if value == nil {
      return "", fmt.Errorf("Diver not found: %s", args[0])
     }
    data := Diver{}
    json.Unmarshal(value, &data)

    level := Level{Levelname:args[1], Org:args[2], Instid: args[3], Status: "Incourse"}
    data.Levels = append(data.Levels, level)

	// JSON  변환
	dataAsBytes, _ := json.Marshal(data)

	err = stub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return string(dataAsBytes), nil
    }

 func addCourse(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	// 해당 아이디가 이전에 있는지 검사
    // GetState => return 값이 nil 아니면 중복 아이디
    value, err := stub.GetState(args[0])
    if err != nil {
      return "", fmt.Errorf("Failed to get diver: %s with error: %s", args[0], err)
     }
    if value == nil {
      return "", fmt.Errorf("Diver not found: %s", args[0])
     }
    data := Diver{}
	json.Unmarshal(value, &data)
	
	lastlevel := len(data.Levels)

	if data.Levels[lastlevel-1].Levelname == args[1] {
		data.Levels[lastlevel-1].Courses = append(data.Levels[lastlevel-1].Courses, args[2])
	}
	
	// JSON  변환
	dataAsBytes, _ := json.Marshal(data)

	err = stub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return string(dataAsBytes), nil
 }

 func addTestResult(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 3 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	// 해당 아이디가 이전에 있는지 검사
    // GetState => return 값이 nil 아니면 중복 아이디
    value, err := stub.GetState(args[0])
    if err != nil {
      return "", fmt.Errorf("Failed to get diver: %s with error: %s", args[0], err)
     }
    if value == nil {
      return "", fmt.Errorf("Diver not found: %s", args[0])
     }
    data := Diver{}
	json.Unmarshal(value, &data)
	
	lastlevel := len(data.Levels)

	if data.Levels[lastlevel-1].Levelname == args[1] {
		data.Levels[lastlevel-1].Status = args[2]
	}

	// JSON  변환
	dataAsBytes, _ := json.Marshal(data)

	err = stub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return string(dataAsBytes), nil
 }
 
 // Get returns the value of the specified asset key
 func getLevel(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	 if len(args) != 1 {
		 return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	 }
 
	 value, err := stub.GetState(args[0])
	 if err != nil {
		 return "", fmt.Errorf("Failed to get diver: %s with error: %s", args[0], err)
	 }
	 if value == nil {
		 return "", fmt.Errorf("diver not found: %s", args[0])
	 }
	 return string(value), nil
 }

 func getHistoryForKey(stub shim.ChaincodeStubInterface, args []string) (string, error) {

	if len(args) < 1 {
	   return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	keyName := args[0]
	// 로그 남기기
	fmt.Println("getHistoryForKey:" + keyName)
 
	resultsIterator, err := stub.GetHistoryForKey(keyName)
	if err != nil {
	   return "", err
	}
	defer resultsIterator.Close()
 
	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")
 
	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
	   response, err := resultsIterator.Next()
	   if err != nil {
		  return "", err
	   }
	   if bArrayMemberAlreadyWritten == true {
		  buffer.WriteString(",")
	   }
	   buffer.WriteString("{\"TxId\":")
	   buffer.WriteString("\"")
	   buffer.WriteString(response.TxId)
	   buffer.WriteString("\"")
 
	   buffer.WriteString(", \"Value\":")
	   if response.IsDelete {
		  buffer.WriteString("null")
	   } else {
		  buffer.WriteString(string(response.Value))
	   }
 
	   buffer.WriteString(", \"Timestamp\":")
	   buffer.WriteString("\"")
	   buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
	   buffer.WriteString("\"")
 
	   buffer.WriteString(", \"IsDelete\":")
	   buffer.WriteString("\"")
	   buffer.WriteString(strconv.FormatBool(response.IsDelete))
	   buffer.WriteString("\"")
 
	   buffer.WriteString("}")
	   bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
 
	// 로그 남기기
	fmt.Println("getHistoryForKey returning:\n" + buffer.String() + "\n")
 
	return (string)(buffer.Bytes()), nil
 }
 
 // main function starts up the chaincode in the container during instantiate
 func main() {
	 if err := shim.Start(new(Dolphins)); err != nil {
		 fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	 }
 }
 