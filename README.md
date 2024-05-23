# IEMS2.0
Blockchain based Excisemanagement System
---------------

In Hyperledger Fabric, handling such scenarios where human errors require corrections while maintaining the integrity and immutability of the blockchain can be approached in several ways. 
Here’s how you can manage this scenario:


**Scenario Description**

**Human Error:** During the manufacturing process, the wrong brand code (x2y) was selected instead of the correct one (xyz) for a batch of 1000 bottles.

**Request:** An approval letter from the commissioner authorizes the correction of the brand code in the database.


**Handling the Scenario in Hyperledger Fabric**


**Recording the Error and Approval:**

First, we record the error and the corresponding approval from the commissioner on the blockchain. This ensures that there is an immutable record of the mistake and the authorization to correct it.


**Implementing Correction Mechanism:**

Develop a chaincode function to handle corrections. This function should only be executable by authorized personnel (e.g., with multi-signature approval or role-based access control).


**Creating an Audit Trail:**


Maintain a detailed audit trail of the correction process, including the original data, the correction request, and the approval.


**Implementation Steps**

**Define the Chaincode for Error Correction:**


This chaincode will handle the recording of errors, requests for correction, and the actual correction process upon approval.


**Example Chaincode:**


package main

import (
    "fmt"
    "encoding/json"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing alcohol manufacturing data
type SmartContract struct {
    contractapi.Contract
}

// BottleData represents the data structure for each batch of bottles
type BottleData struct {
    BatchID     string `json:"batchID"`
    BrandCode   string `json:"brandCode"`
    Quantity    int    `json:"quantity"`
    Error       bool   `json:"error"`
    ErrorDetails string `json:"errorDetails"`
}

// ErrorCorrectionRequest represents a request to correct an error
type ErrorCorrectionRequest struct {
    BatchID        string `json:"batchID"`
    IncorrectBrand string `json:"incorrectBrand"`
    CorrectBrand   string `json:"correctBrand"`
    ApprovedBy     string `json:"approvedBy"`
    ApprovalLetter string `json:"approvalLetter"`
}

// InitLedger adds a base set of data to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
    bottles := []BottleData{
        {BatchID: "1", BrandCode: "xyz", Quantity: 1000, Error: false, ErrorDetails: ""},
    }

    for _, bottle := range bottles {
        bottleJSON, err := json.Marshal(bottle)
        if err != nil {
            return err
        }

        err = ctx.GetStub().PutState(bottle.BatchID, bottleJSON)
        if err != nil {
            return fmt.Errorf("failed to put to world state. %v", err)
        }
    }

    return nil
}

// RecordError records an error in a batch
func (s *SmartContract) RecordError(ctx contractapi.TransactionContextInterface, batchID string, errorDetails string) error {
    bottleJSON, err := ctx.GetStub().GetState(batchID)
    if err != nil {
        return fmt.Errorf("failed to read from world state. %v", err)
    }
    if bottleJSON == nil {
        return fmt.Errorf("batch %s does not exist", batchID)
    }

    var bottle BottleData
    err = json.Unmarshal(bottleJSON, &bottle)
    if err != nil {
        return err
    }

    bottle.Error = true
    bottle.ErrorDetails = errorDetails

    updatedBottleJSON, err := json.Marshal(bottle)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(batchID, updatedBottleJSON)
}

// RequestCorrection records a request for error correction
func (s *SmartContract) RequestCorrection(ctx contractapi.TransactionContextInterface, batchID string, incorrectBrand string, correctBrand string, approvedBy string, approvalLetter string) error {
    correctionRequest := ErrorCorrectionRequest{
        BatchID:        batchID,
        IncorrectBrand: incorrectBrand,
        CorrectBrand:   correctBrand,
        ApprovedBy:     approvedBy,
        ApprovalLetter: approvalLetter,
    }

    correctionRequestJSON, err := json.Marshal(correctionRequest)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState("CORRECTION_"+batchID, correctionRequestJSON)
}

// CorrectError corrects the error in a batch upon approval
func (s *SmartContract) CorrectError(ctx contractapi.TransactionContextInterface, batchID string) error {
    correctionRequestJSON, err := ctx.GetStub().GetState("CORRECTION_" + batchID)
    if err != nil {
        return fmt.Errorf("failed to read from world state. %v", err)
    }
    if correctionRequestJSON == nil {
        return fmt.Errorf("correction request for batch %s does not exist", batchID)
    }

    var correctionRequest ErrorCorrectionRequest
    err = json.Unmarshal(correctionRequestJSON, &correctionRequest)
    if err != nil {
        return err
    }

    bottleJSON, err := ctx.GetStub().GetState(batchID)
    if err != nil {
        return fmt.Errorf("failed to read from world state. %v", err)
    }
    if bottleJSON == nil {
        return fmt.Errorf("batch %s does not exist", batchID)
    }

    var bottle BottleData
    err = json.Unmarshal(bottleJSON, &bottle)
    if err != nil {
        return err
    }

    // Update the brand code
    bottle.BrandCode = correctionRequest.CorrectBrand
    bottle.Error = false
    bottle.ErrorDetails = ""

    updatedBottleJSON, err := json.Marshal(bottle)
    if err != nil {
        return err
    }

    return ctx.GetStub().PutState(batchID, updatedBottleJSON)
}

func main() {
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    if err != nil {
        fmt.Printf("Error create alcohol manufacturing chaincode: %s", err.Error())
        return
    }

    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting alcohol manufacturing chaincode: %s", err.Error())
    }
}



**Steps to Implement**


**Record the Error:**

When the error is discovered, the RecordError function is called to update the batch record with error details.


**Request Correction:**

The RequestCorrection function records the request for correction along with the approval from the commissioner.

**Approval and Correction:**

After verifying the approval, the CorrectError function is called to update the incorrect brand code with the correct one.


**Audit Trail:**

All steps (recording the error, requesting correction, and performing the correction) are logged on the blockchain, creating an immutable audit trail.


**Advantages of Using Hyperledger Fabric**

**Immutability:** Even though corrections can be made, all actions are logged immutably on the blockchain, preserving a history of changes.


**Permissioned Access:** Only authorized individuals can request and approve corrections, ensuring security and compliance.


**Traceability:** Detailed audit trails help in maintaining transparency and accountability.


**Automated Workflows:** Smart contracts automate the validation and correction processes, reducing the risk of further human error.


This approach ensures that the integrity of the blockchain is maintained while allowing for necessary corrections under controlled and authorized conditions.
