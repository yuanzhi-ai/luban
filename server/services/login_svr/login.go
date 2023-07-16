package main

import (
	"context"
	"fmt"

	"github.com/yuanzhi-ai/luban/go_proto/verify_proto"
	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/data"
	"github.com/yuanzhi-ai/luban/server/log"
	"github.com/yuanzhi-ai/luban/server/sms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getMachineVerify(ctx context.Context) (string, string, int32, error) {
	verifyReq := &verify_proto.GetVerCodeReq{}
	// 这里走grpc的客户端
	grpcConn, err := grpc.Dial("172.31.66.86:3360", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Dial chat svr err:%v", err)
		return "", "", comm.NetWorkErr, err
	}
	client := verify_proto.NewVerifyClient(grpcConn)
	defer grpcConn.Close()
	rsp, err := client.GetVerCode(ctx, verifyReq)
	if err != nil {
		rspErr := fmt.Errorf("RPC GetVerCode fial err:%v", err)
		log.Errorf("%v", rspErr)
		return "", "", comm.GetVerCodeErr, rspErr
	}
	return rsp.CodeId, rsp.Base64Img, comm.SuccessCode, nil
}

// 验证用户输入的验证码
func sendMachineVerifyResult(ctx context.Context, capId string, userAns string) (int32, error) {
	verifyReq := &verify_proto.VerifyCodeReq{CodeId: capId, Ans: userAns}
	grpcConn, err := grpc.Dial("172.31.66.86:3360", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Dial chat svr err:%v", err)
		return comm.NetWorkErr, err
	}
	client := verify_proto.NewVerifyClient(grpcConn)
	defer grpcConn.Close()
	rsp, err := client.VerifyCode(ctx, verifyReq)
	if err != nil {
		rspErr := fmt.Errorf("RPC VerifyCode fial err:%v", err)
		log.Errorf("%v", rspErr)
		return comm.VerifyCodeErr, rspErr
	}
	return rsp.RetCode, nil
}

const (
	// 手机验证码长度
	codeLen = 6
)

// 向手机发送一个验证码
func sendTextVerCode(ctx context.Context, phone string) (int32, error) {
	// 判读手机号是否合法
	if !comm.IsPhoneLegal(phone) {
		return comm.InputErr, fmt.Errorf("input phone:%v legal", phone)
	}
	code := comm.GetRandDigitStr(codeLen)
	err := sms.SendMsg(code, phone)
	if err != nil {
		log.Errorf("send vercode fial phone:%v code:%v err:%v", phone, code, err)
		return comm.SmsErr, err
	}
	// 向redis写入
	err = data.SetPhoneCaptcha(ctx, phone, code)
	if err != nil {
		log.Errorf("set phone code to redis fial. phone:%v code:%v err:%v", phone, code, err)
		return comm.RedisErr, err
	}
	return comm.SuccessCode, nil
}

// 验证用户输入的手机验证码
func verifyPhoneCode(ctx context.Context, phone string, inCode string) (int32, error) {
	// 判读验证码是否合法
	if !comm.IsPhoneCodeLegal(inCode, codeLen) {
		return comm.InputErr, fmt.Errorf("input code:%v legal", inCode)
	}
	// 判读手机号是否合法
	if !comm.IsPhoneLegal(phone) {
		return comm.InputErr, fmt.Errorf("input phone:%v legal", phone)
	}
	// 查询redis
	code, err := data.GetPhoneCaptcha(ctx, phone)
	if err != nil {
		log.Errorf("get phone captcha fail. phone:%v err:%v", phone, err)
		return comm.QueryPhoneCodeErr, err
	}
	if code == "" {
		return comm.CodeExpiration, nil
	}
	if inCode != code {
		return comm.PhoneCodeWrong, nil
	}
	// 走到这里说明code验证没问题，需要删除用过的code
	data.DelPhoneCaptcha(ctx, phone)
	return comm.SuccessCode, nil
}

// register 用户注册
func register(ctx context.Context, phone string, code string, pswd string) (int32, error) {
	retCode, err := verifyPhoneCode(ctx, phone, code)
	if err != nil || retCode != comm.SuccessCode {
		return retCode, err
	}
	// 向数据库插入用户记录
	return comm.SuccessCode, nil
}

// passwordLogin 手机号+密码登录
func passwordLogin(ctx context.Context, phone string, al string) (int32, error) {
	return comm.SuccessCode, nil
}
