package model

import "time"

type Product struct {
	ID                int32     `json:"id"`
	Name              string    `json:"name"`
	ProductType       int32     `gorm:"default:1" json:"productType"`
	CategoryId        int32     `json:"categoryId"`
	StartingPrice     float32   `json:"startingPrice"`
	TotalStock        int32     `gorm:"default:1234" json:"totalStock"`
	MainPicture       string    `gorm:"default:1" json:"mainPicture"`
	RemoteAreaPostage float32   `json:"remoteAreaPostage"`
	SingleBuyLimit    int32     `json:"singleBuyLimit"`
	IsEnable          bool      `json:"isEnable"`
	Remark            string    `gorm:"default:1" json:"remark"`
	CreateUser        int32     `gorm:"default:1" json:"createUser"`
	CreateTime        time.Time `json:"createTime"`
	UpdateUser        int32     `json:"updateUser"`
	UpdateTime        time.Time `json:"updateTime"`
	IsDeleted         bool      `json:"isDeleted"`
}

func (table *Product) TableName() string {
	return "product"
}
