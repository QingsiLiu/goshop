package handler

import (
    "common"
    "context"
    "fmt"
    "strconv"
    "time"
    "user-ser/domain/model"
    "user-ser/domain/service"
    "user-ser/proto"
)

type UserHandler struct {
    UserDataService service.IUserDataService
}

// Login 注册登录
func (u *UserHandler) Login(ctx context.Context, request *proto.LoginRequest, resp *proto.LoginResp) error {
    userInfo, err := u.UserDataService.Login(request.GetClientId(), request.GetSystemId(), request.GetPhone(), request.GetVerificationCode())
    if err != nil {
        return err
    }
    fmt.Println(">>>>>>>>>>>>> login success :", userInfo)

    userForResp(userInfo, resp)
    u.UserDataService.SetUserToken(strconv.Itoa(int(userInfo.ID)), []byte(resp.Token), time.Duration(1)*time.Hour)
    return nil
}

func userForResp(user *model.User, resp *proto.LoginResp) *proto.LoginResp {
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

func (u *UserHandler) GetUserToken(ctx context.Context, req *proto.TokenReq, resp *proto.TokenResp) error {
    res, _ := u.UserDataService.GetUserToken(req.GetUuid())
    if res != "" {
        resp.IsLogin = true
        resp.Token = res
        // 续命
        uuid := common.ToInt(req.Uuid)
        u.UserDataService.SetUserToken(strconv.Itoa(uuid), []byte(res), time.Duration(1)*time.Hour)
        fmt.Println(">>>>>>>>>>>>> GetUserToken success :", res)
    } else {
        resp.IsLogin = false
        resp.Token = ""
    }
    return nil
}
