package model

import (
	"time"
	"errors"
	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type PixKey struct {
	Base               `valid:"required"`
	Kind      string   `json:"kind" gorm:"type:varchar(20);not null" valid:"notnull"`
	Key       string   `json:"key" gorm:"type:varchar(255);not null" valid:"notnull"`
	Account   *Account `valid:"-"`
	AccountID string   `gorm:"column:account_id;type:uuid;not null" valid:"-"`
	Status    string   `json:"status" gorm:"type:varchar(20);not null" valid:"notnull"`
}

func (pixKey *PixKey) isValid() error {
	_, err := govalidator.ValidateStruct(pixKey)

	if pixKey.Kind != "email" && pixKey.Kind != "cpf" {
		return errors.New("invalid type of key")
	}

	if pixKey.Status != "active" && pixKey.Status != "inactive" {
		return errors.New("invalid status")
	}

	if err != nil {
		return err
	}
	return nil
}

func NewPixKey(kind string, account *Account, key string) (*PixKey, error) {
	pixKey := PixKey {
		Kind:			kind,
		Key:			key,
		Account:	account,
		Status:		"active",
	}

	pixKey.ID = uuid.NewV4().String()
	pixKey.CreatedAt = time.Now();

	err := pixKey.isValid()
	if err != nil {
		return nil, err
	}

	return &pixKey, nil
}