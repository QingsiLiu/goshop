package service

import (
    "time"
    "user-ser/domain/model"
    "user-ser/domain/repository"
)

type IUserDataService interface {
    Login(clientId, systemId int32, phone, verificationCode string) (*model.User, error)
    SetUserToken(key string, val []byte, timeTTL time.Duration)
    GetUserToken(key string) (string, error)
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

func (u *UserDataService) SetUserToken(key string, val []byte, timeTTL time.Duration) {
    u.userRepo.SetUserToken(key, val, timeTTL)
}

func (u *UserDataService) GetUserToken(key string) (string, error) {
    return u.userRepo.GetUserToken(key)
}
