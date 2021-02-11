package usecase

import (
	"errors"
	"github.com/MrChampz/codepix/domain/model"
	"github.com/MrChampz/codepix/domain/repository"
)

type PixUseCase struct {
	PixRepository repository.PixRepositoryInterface
}

func (useCase *PixUseCase) RegisterKey(
	key string,
	kind string,
	accountId string,
) (*model.PixKey, error) {
	account, err := useCase.PixRepository.FindAccount(accountId)
	if err != nil {
		return nil, err
	}

	pixKey, err := model.NewPixKey(kind, account, key)
	if err != nil {
		return nil, err
	}

	err = useCase.PixRepository.RegisterKey(pixKey)
	if err != nil {
		return nil, errors.New("unable to create this key now")
	}

	return pixKey, nil
}

func (useCase *PixUseCase) FindKey(key string, kind string) (*model.PixKey, error) {
	return useCase.PixRepository.FindKeyByKind(key, kind)
}