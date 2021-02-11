package factory

import (
	"github.com/MrChampz/codepix/application/usecase"
	"github.com/MrChampz/codepix/infrastructure/repository"
	"github.com/jinzhu/gorm"
)

func TransactionUseCaseFactory(database *gorm.DB) usecase.TransactionUseCase {
	pixRepository := repository.PixRepositoryDb { Db: database }
	transactionRepository := repository.TransactionRepositoryDb { Db: database }

	transactionUseCase := usecase.TransactionUseCase {
		TransactionRepository:	&transactionRepository,
		PixRepository:					pixRepository,
	}

	return transactionUseCase
}