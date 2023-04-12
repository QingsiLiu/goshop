package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"user-ser/common"
	"user-ser/proto"
)

func UserLogin(c *gin.Context) {
	// 获取远程服务的客户端(获取服务)
	client := getClient()

	// 获取页面参数，拼接请求信息
	clientId, _ := strconv.Atoi(c.Request.FormValue("clientId"))
	phone := c.Request.FormValue("phone")
	systemId, _ := strconv.Atoi(c.Request.FormValue("systemId"))
	verificationCode := c.Request.FormValue("verificationCode")

	req := &proto.LoginRequest{
		ClientId:         int32(clientId),
		Phone:            phone,
		SystemId:         int32(systemId),
		VerificationCode: verificationCode,
	}

	// 远程调用服务
	resp, err := client.Login(context.TODO(), req)
	// 根据响应输出
	if err != nil {
		log.Println(err.Error())
		common.RespFail(c.Writer, "登录失败", resp)
		return
	}

	common.RespOK(c.Writer, "登录成功", resp)
}
