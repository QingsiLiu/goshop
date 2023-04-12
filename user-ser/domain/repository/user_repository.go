package repository

import (
	"errors"
	"gorm.io/gorm"
	"user-ser/domain/model"
)

type IUserRepository interface {
	Login(clientId, systemId int32, phone, verificationCode string) (*model.User, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		mysqlDB: db,
	}
}

type UserRepository struct {
	mysqlDB *gorm.DB
}

func (u *UserRepository) Login(clientId, systemId int32, phone, verificationCode string) (user *model.User, err error) {
	user = &model.User{}
	if clientId != 0 && systemId != 0 && verificationCode == "6666" {
		u.mysqlDB.Where("phone = ?", phone).Find(user)
		// 未找到就注册一个
		if user.ID == 0 {
			user.Phone = phone
			u.mysqlDB.Create(&user)
		}
		return user, nil
	} else {
		return user, errors.New("参数不匹配")
	}
}
