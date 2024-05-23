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
