package service

import (
	"user-ser/domain/model"
	"user-ser/domain/repository"
)

type IUserDataService interface {
	Login(clientId, systemId int32, phone, verificationCode string) (*model.User, error)
}

func NewUserDataService(repo repository.IUserRepository) IUserDataService {
	return &UserDataService{
		userRepo: repo,
	}
}

type UserDataService struct {
	userRepo repository.IUserRepository
}

func (u *UserDataService) Login(clientId, systemId int32, phone, verificationCode string) (*model.User, error) {

	return u.userRepo.Login(clientId, systemId, phone, verificationCode)
}
