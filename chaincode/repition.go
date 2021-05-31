package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Data
type SmartContract struct {
	contractapi.Contract
}

// Data describes basic details of what makes up a simple repition data
type Data struct {
	ID              string `json:"id"`
	TaskId          int    `json:"task_id"`
	CompanyCode     int    `json:"company_code"`
	LetterNum       string `json:"letter_num"`
	PredictValue    string `json:"predict_value"`
	PredictDivation string `json:"predict_divation"`
	Fvalue          string `json:"f_value"`
}

// InitLedger adds a base set of datas to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	datas := []Data{
		{ID: "1", TaskId: 122, CompanyCode: 2, LetterNum: "ln1", PredictValue: "pv1", PredictDivation: "pd1", Fvalue: "fv1"},
		{ID: "2", TaskId: 122, CompanyCode: 3, LetterNum: "ln2", PredictValue: "pv2", PredictDivation: "pd2", Fvalue: "fv2"},
	}

	for _, data := range datas {
		dataJSON, err := json.Marshal(data)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(data.ID, dataJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// Createdata issues a new data to the world state with given details.
func (s *SmartContract) CreateData(ctx contractapi.TransactionContextInterface, id string, task_id int, company_code int, letter_num string, predict_value string, predict_divation string, f_value string) error {
	exists, err := s.DataExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the data %s already exists", id)
	}

	data := Data{
		ID:              id,
		TaskId:          task_id,
		CompanyCode:     company_code,
		LetterNum:       letter_num,
		PredictValue:    predict_value,
		PredictDivation: predict_divation,
		Fvalue:          f_value,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, dataJSON)
}

// ReadData returns the data stored in the world state with given id.
func (s *SmartContract) ReadData(ctx contractapi.TransactionContextInterface, id string) (*Data, error) {
	dataJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if dataJSON == nil {
		return nil, fmt.Errorf("the data %s does not exist", id)
	}

	var data Data
	err = json.Unmarshal(dataJSON, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// UpdateData updates an existing data in the world state with provided parameters.
func (s *SmartContract) UpdateData(ctx contractapi.TransactionContextInterface, id string, task_id int, company_code int, letter_num string, predict_value string, predict_divation string, f_value string) error {
	exists, err := s.DataExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the data %s does not exist", id)
	}

	// overwriting original data with new data
	data := Data{
		ID:              id,
		TaskId:          task_id,
		CompanyCode:     company_code,
		LetterNum:       letter_num,
		PredictValue:    predict_value,
		PredictDivation: predict_divation,
		Fvalue:          f_value,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, dataJSON)
}

// DeleteData deletes an given data from the world state.
func (s *SmartContract) DeleteData(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.DataExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the data %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// DataExists returns true when data with given id exists in world state
func (s *SmartContract) DataExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	dataJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return dataJSON != nil, nil
}

// GetAllDatas returns all datas found in world state
func (s *SmartContract) GetAllDatas(ctx contractapi.TransactionContextInterface) ([]*Data, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all datas in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var datas []*Data
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var data Data
		err = json.Unmarshal(queryResponse.Value, &data)
		if err != nil {
			return nil, err
		}
		datas = append(datas, &data)
	}

	return datas, nil
}
