package repository

import "github.com/MrChampz/codepix/domain/model"

type TransactionRepositoryInterface interface {
	Register(transaction *model.Transaction) error
	Save(transaction *model.Transaction) error
	Find(id string) (*model.Transaction, error)
}