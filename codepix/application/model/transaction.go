package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type Transaction struct {
	ID	 						string		`json:"id" validate:"required,uuid4"`
	AccountID 			string		`json:"accountId" validate:"required,uuid4"`
	Amount 					float64		`json:"amount" validate:"required,numeric"`
	PixKeyTo 				string		`json:"pixKeyTo" validate:"required"`
	PixKeyToKind		string		`json:"pixKeyToKind" valid:"required"`
	Description 		string		`json:"description" validate:"required"`
	Status 					string		`json:"status" validate:"-"`
	Error						string		`json:"error"`
}

func NewTransaction() *Transaction {
	return &Transaction {}
}

func (transaction *Transaction) isValid() error {
	v := validator.New()
	err := v.Struct(transaction)
	if err != nil {
		fmt.Errorf("Error during Transaction validation: %s", err.Error())
		return err
	}
	return nil
}

func (transaction *Transaction) ParseJson(data []byte) error {
	err := json.Unmarshal(data, transaction)
	if err != nil {
		return err
	}
	return transaction.isValid()
}

func (transaction *Transaction) ToJson() ([]byte, error) {
	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(transaction)
	if err != nil {
		return nil, err
	}

	return result, nil
}