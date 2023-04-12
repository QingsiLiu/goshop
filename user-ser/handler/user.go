package handler

import (
	"context"
	"fmt"
	"time"
	"user-ser/common"
	"user-ser/domain/model"
	"user-ser/domain/service"
	"user-ser/proto"
)

type User struct {
	UserDataService service.IUserDataService
}

// Login 注册登录
func (u *User) Login(ctx context.Context, request *proto.LoginRequest, resp *proto.LoginResp) error {
	userInfo, err := u.UserDataService.Login(request.ClientId, request.SystemId, request.Phone, request.VerificationCode)
	if err != nil {
		return err
	}
	fmt.Println(">>>>>>>>>>>>> login success :", userInfo)
	UserForResp(userInfo, resp)
	return nil
}

func UserForResp(user *model.User, resp *proto.LoginResp) *proto.LoginResp {
	timeStr := fmt.Sprintf("%d", time.Now().Unix())
	resp.Token = common.Md5Encode(timeStr)

	tp, _ := time.ParseDuration("1h")
	tokenExpireTime := time.Now().Add(tp)
	expireTimeStr := tokenExpireTime.Format("2006-01-02 15:04:05")

	resp.User = &proto.User{
		Avatar:          user.Avatar,
		ClientId:        user.ClientId,
		EmployeeId:      1,
		Nickname:        user.Nickname,
		Phone:           user.Phone,
		SessionId:       resp.Token,
		Token:           resp.Token,
		TokenExpireTime: expireTimeStr,
		UnionId:         user.UnionId,
		Id:              user.ID,
	}

	return resp
}
