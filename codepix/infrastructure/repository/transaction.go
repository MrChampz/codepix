package repository

import (
	"fmt"
	"github.com/MrChampz/codepix/domain/model"
	"github.com/jinzhu/gorm"
)

type TransactionRepositoryDb struct {
	Db *gorm.DB
}

func (repository *TransactionRepositoryDb) Register(transaction *model.Transaction) error {
	return repository.Db.Create(transaction).Error
}

func (repository *TransactionRepositoryDb) Save(transaction *model.Transaction) error {
	return repository.Db.Save(transaction).Error
}

func (repository *TransactionRepositoryDb) Find(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	repository.Db.Preload("AccountFrom.Bank").First(&transaction, "id = ?", id)
	if transaction.ID == "" {
		return nil, fmt.Errorf("no transaction was found")
	}
	return &transaction, nil
}