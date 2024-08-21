package usecase

import (
	"taskTest/internal/entity"
	"taskTest/internal/storage"
)

type useCase struct {
	storage storage.RepositoryDB
}

func NewUseCase(storage storage.RepositoryDB) *useCase {
	return &useCase{storage: storage}
}
func (uc *useCase) CreateUser(user *entity.User) error {
	return uc.storage.CreateUser(user)
}

func (uc *useCase) GetUserByID(id int) (*entity.User, error) {
	return uc.storage.GetUserByID(id)
}

func (uc *useCase) UpdateUser(user *entity.User) error {
	return uc.storage.UpdateUser(user)
}

func (uc *useCase) DeleteUser(id int) error {
	return uc.storage.DeleteUser(id)
}
