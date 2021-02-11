package repository

import (
	"fmt"
	"github.com/MrChampz/codepix/domain/model"
	"github.com/jinzhu/gorm"
)

type PixRepositoryDb struct {
	Db *gorm.DB
}

func (repository PixRepositoryDb) AddBank(bank *model.Bank) error {
	return repository.Db.Create(bank).Error
}

func (repository PixRepositoryDb) AddAccount(account *model.Account) error {
	return repository.Db.Create(account).Error
}

func (repository PixRepositoryDb) RegisterKey(key *model.PixKey) error {
	return repository.Db.Create(key).Error
}

func (repository PixRepositoryDb) FindBank(id string) (*model.Bank, error) {
	var bank model.Bank
	repository.Db.Preload("Account").First(&bank, "id = ?", id)
	if bank.ID == "" {
		return nil, fmt.Errorf("no bank found")
	}
	return &bank, nil
}

func (repository PixRepositoryDb) FindAccount(id string) (*model.Account, error) {
	var account model.Account
	repository.Db.Preload("Bank").First(&account, "id = ?", id)
	if account.ID == "" {
		return nil, fmt.Errorf("no account found")
	}
	return &account, nil
}

func (repository PixRepositoryDb) FindKeyByKind(key string, kind string) (*model.PixKey, error) {
	var pixKey model.PixKey
	repository.Db.Preload("Account.Bank").First(&pixKey, "kind = ? and key = ?", kind, key)
	if pixKey.ID == "" {
		return nil, fmt.Errorf("no key was found")
	}
	return &pixKey, nil
}