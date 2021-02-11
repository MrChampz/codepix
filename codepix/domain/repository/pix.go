package repository

import "github.com/MrChampz/codepix/domain/model"

type PixRepositoryInterface interface {
	AddBank(bank *model.Bank) error
	FindBank(id string) (*model.Bank, error)
	AddAccount(account *model.Account) error
	FindAccount(id string) (*model.Account, error)
	RegisterKey(key *model.PixKey) error
	FindKeyByKind(key string, kind string) (*model.PixKey, error)
}