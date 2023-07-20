package main

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	// 判断手机号是否合法
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
	userInfo, err := generatorUserRegister(phone, pswd)
	if err != nil {
		return comm.GeneratorUserInfoErr, err
	}
	err = insterUserRegristerInfo(ctx, userInfo)
	if err != nil {
		return comm.InsertUserRegisterInfoErr, err
	}
	return comm.SuccessCode, nil
}

// passwordLogin 手机号+密码登录
func passwordLogin(ctx context.Context, phone string, a1 string) (int32, error) {
	//从数据库拿用户的DB_A1
	dbS2, err := getUserS2(ctx, phone)
	if err != nil {
		return comm.QueryS2Err, fmt.Errorf("getUserS2 fail. err:%v", err)
	}
	// 对dbs2做16进制解码
	dbKey, err := hex.DecodeString(dbS2)
	if err != nil {
		return comm.HexDBS2Err, fmt.Errorf("hex Decode fail. dbS2:%v err:%v", dbS2, err)
	}
	// 对前端传来的a1做base64解码
	uA1, err := base64.StdEncoding.DecodeString(a1)
	if err != nil {
		return comm.DecodeA1Err, err
	}
	// 尝试使用dbS2对uA1解码
	decA1 := comm.AesDecryptCBC(uA1, dbKey)
	a1Data, err := getA1ValueFromData(string(decA1))
	if err != nil {
		return comm.DercyA1Err, err
	}
	// 比较解码的phone与请求的phone是否一致
	if phone != a1Data.phone {
		return comm.LoginErr, fmt.Errorf("check login fail, diff phone. req phone:%v decrypt phone:%v", phone, a1Data.phone)
	}
	// 比较时间戳
	now := time.Now().Unix()
	tsGap := now - a1Data.ts
	if tsGap > 5*60 || tsGap < -5*60 {
		return comm.LoginErr, fmt.Errorf("check loging ts fail. svr ts:%v client:%v tsGap:%v", now, a1Data.ts, tsGap)
	}
	// 比较重新计算的s2 与db s2
	clinetS2 := comm.Md5Encode(a1Data.phone + a1Data.md5Pswd)
	if clinetS2 != dbS2 {
		return comm.LoginErr, fmt.Errorf("check s2 fail. clients2:%v dbs2:%v", clinetS2, dbS2)
	}
	// 更新时间戳 柔性
	err = updateLoginTs(ctx, phone)
	if err != nil {
		log.Errorf("update login ts err:%v", err)
	}
	return comm.SuccessCode, nil
}

type a1Data struct {
	phone string // 手机号

	md5Pswd string // md5(password)

	ts int64 // 登录的前端时间戳

	randomKey string // 随机key
}

// 获取解码后的数据
func getA1ValueFromData(uData string) (*a1Data, error) {

	datas := strings.Split(uData, ";")
	if len(datas) != 4 {
		return nil, fmt.Errorf("decode A1Data fail. UData:%v", uData)
	}
	ts, err := strconv.ParseInt(datas[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("decode A1Data fail, ts is wrong. ts:%v", datas[2])
	}
	uA1 := &a1Data{
		phone:     datas[0],
		md5Pswd:   datas[1],
		ts:        ts,
		randomKey: datas[3],
	}
	return uA1, nil
}
