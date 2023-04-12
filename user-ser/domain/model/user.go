package model

import "time"

type User struct {
	ID            int32
	Avatar        string `gorm:"default:'https://msb-edu-dev.oss-cn-beijing.aliyuncs.com/default-headimg.png'"`
	ClientId      int32  `gorm:"default:1"`
	Nickname      string `gorm:"default:'随机名称'"`
	Phone         string
	Password      string `gorm:"default:'1234'"`
	SystemId      string `gorm:"default:1"`
	LastLoginTime time.Time
	CreateTime    time.Time
	IsDeleted     int32  `gorm:"default:0"`
	UnionId       string `gorm:"default:'1'"`
}

func (table *User) TableName() string {
	return "user"
}
