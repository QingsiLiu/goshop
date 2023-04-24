package repository

import (
    "common"
    "errors"
    "fmt"
    "github.com/go-redis/redis/v8"
    "gorm.io/gorm"
    "time"
    "user-ser/domain/model"
)

type IUserRepository interface {
    Login(clientId, systemId int32, phone, verificationCode string) (*model.User, error)
    SetUserToken(key string, val []byte, timeTTL time.Duration)
    GetUserToken(key string) (string, error)
}

func NewUserRepository(db *gorm.DB, red *redis.Client) IUserRepository {
    return &UserRepository{
        mysqlDB: db,
        red:     red,
    }
}

type UserRepository struct {
    mysqlDB *gorm.DB
    red     *redis.Client
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

func (u *UserRepository) SetUserToken(key string, val []byte, timeTTL time.Duration) {

    intKey := common.ToInt(key)
    binKey := common.ConverToBinary(intKey)
    fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>", binKey)

    common.SetUserToken(u.red, binKey, val, timeTTL)
}

func (u *UserRepository) GetUserToken(key string) (string, error) {
    return common.GetUserToken(u.red, key)
}
