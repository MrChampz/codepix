package usecase

import (
	"errors"
	"github.com/MrChampz/codepix/domain/model"
	"github.com/MrChampz/codepix/domain/repository"
	"log"
)

type TransactionUseCase struct {
	TransactionRepository repository.TransactionRepositoryInterface
	PixRepository					repository.PixRepositoryInterface
}

func (useCase *TransactionUseCase) Register(
	id string,
	accountId string,
	amount float64,
	pixKeyTo string,
	pixKeyToKind string,
	description string,
) (*model.Transaction, error) {

	account, err := useCase.PixRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := useCase.PixRepository.FindKeyByKind(pixKeyTo, pixKeyToKind)
	if err != nil {
		return nil, err
	}

	transaction, err := model.NewTransaction(id, account, amount, pixKey, description)
	if err != nil {
		return nil, err
	}

	err = useCase.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, errors.New("unable to process this transaction")
	}

	return transaction, nil
}

func (useCase *TransactionUseCase) Confirm(
	transactionId string,
) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.Find(transactionId)
	if err != nil {
		log.Println("Transaction not found", transactionId)
		return nil, err
	}

	transaction.Confirm()

	err = useCase.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (useCase *TransactionUseCase) Complete(
	transactionId string,
) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.Find(transactionId)
	if err != nil {
		log.Println("Transaction not found", transactionId)
		return nil, err
	}

	transaction.Complete()

	err = useCase.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (useCase *TransactionUseCase) Error(
	transactionId string,
	reason string,
) (*model.Transaction, error) {
	transaction, err := useCase.TransactionRepository.Find(transactionId)
	if err != nil {
		log.Println("Transaction not found", transactionId)
		return nil, err
	}

	transaction.Cancel(reason)

	err = useCase.TransactionRepository.Save(transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}