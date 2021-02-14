package model

import (
	"time"
	"errors"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending string = "pending"
	TransactionConfirmed string = "confirmed"
	TransactionCompleted string = "completed"
	TransactionError string = "error"
)

type Transaction struct {
	Base 												`valid:"required"`
	AccountFrom 			*Account	`valid:"-"`
	AccountFromID			string		`gorm:"column:account_from_id;type:uuid;" valid:"notnull"`
	Amount 						float64		`json:"amount" gorm:"type:float" valid:"notnull"`
	PixKeyTo 					*PixKey		`valid:"-"`
	PixKeyToID				string		`gorm:"column:pix_key_to_id;type:uuid;" valid:"notnull"`
	Status 						string		`json:"status" gorm:"type:varchar(20)" valid:"notnull"`
	Description 			string		`json:"description" gorm:"type:varchar(255)" valid:"-"`
	CancelDescription string		`json:"cancel_description" gorm:"type:varchar(255)" valid:"-"`
}

func (transaction *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(transaction)

	if transaction.Amount <= 0 {
		return errors.New("the amount must be greater than 0")
	}

	if transaction.Status != TransactionPending &&
	   transaction.Status != TransactionCompleted &&
	   transaction.Status != TransactionError {
		return errors.New("invalid status for the transaction")
	}

	if transaction.PixKeyTo != nil &&
	   transaction.PixKeyTo.AccountID == transaction.AccountFrom.ID {
		return errors.New("the source and destination account cannot be the same")
	}

	if err != nil {
		return err
	}
	return nil
}

func NewTransaction(
	id string,
	accountFrom *Account,
	amount float64,
	pixKeyTo *PixKey,
	description string,
) (*Transaction, error) {
	transaction := Transaction {
		AccountFrom: accountFrom,
		AccountFromID: accountFrom.ID,
		Amount: amount,
		PixKeyTo: pixKeyTo,
		PixKeyToID: pixKeyTo.ID,
		Status: TransactionPending,
		Description: description,
	}

	if id == "" {
 		transaction.ID = uuid.NewV4().String()
	} else {
		transaction.ID = id
	}

	transaction.CreatedAt = time.Now()

	err := transaction.isValid()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (transaction *Transaction) Complete() error {
	transaction.Status = TransactionCompleted
	transaction.UpdatedAt = time.Now()
	return transaction.isValid()
}

func (transaction *Transaction) Confirm() error {
	transaction.Status = TransactionConfirmed
	transaction.UpdatedAt = time.Now()
	return transaction.isValid()
}

func (transaction *Transaction) Cancel(description string) error {
	transaction.Status = TransactionError
	transaction.CancelDescription = description
	transaction.UpdatedAt = time.Now()
	return transaction.isValid()
}