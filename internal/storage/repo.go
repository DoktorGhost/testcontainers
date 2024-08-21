package storage

import (
	"taskTest/internal/entity"
)

type RepositoryDB interface {
	CreateUser(user *entity.User) error
	GetUserByID(id int) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id int) error
}
