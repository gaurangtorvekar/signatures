package main

import (
	"fmt"
	// "strconv"
	"encoding/json"
	s "strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type Signature struct {
	Emails  string `json:"email"`
	PdfHash string `json:"pdfhash"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("########### example_cc Init ###########")
	_, args := stub.GetFunctionAndParameters()
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	// Write the state to the ledger
	err = stub.PutState("initialize_var", []byte(args[0]))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)

}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call")
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("invoke is running ")
	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0] == "query" {
		// queries an entity state
		return t.query(stub, args)
	}
	if args[0] == "write" {
		// Adds a new signature to the state
		return t.write(stub, args)
	}
	return shim.Error("Unknown action, check the first argument, must be one of 'delete', 'query', or 'move'")
}

func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var email, hash string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2!")
	}

	email = args[1] //rename for funsies
	hash = args[0]

	// Check if PDF has already been signed by this user
	sigAsBytes, err := stub.GetState(hash)
	if err != nil {
		return shim.Error("Failed to get signature hash")
	}
	res := Signature{}
	json.Unmarshal(sigAsBytes, &res)
	old_emails_str := res.Emails
	old_emails := s.Split(old_emails_str, ",")
	new_emails := []string{}
	new_emails = append(new_emails, email)
	if res.PdfHash == hash {
		for _, v := range old_emails {
			if v != email {
				new_emails = append(new_emails, v)
			}
		}
	}
	new_emails_string := s.Join(new_emails, ",")
	// build the signatures json string manually
	str := `{"email": "` + new_emails_string + `", "pdfhash": "` + hash + `"}`
	err = stub.PutState(hash, []byte(str))
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var key string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valAsbytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
