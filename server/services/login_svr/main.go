package main

import (
	"context"
	"fmt"
	"net"

	"github.com/yuanzhi-ai/luban/go_proto/login_proto"
	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/log"
	"google.golang.org/grpc"
)

type server struct {
	login_proto.UnimplementedLoginServer
}

const (
	PORT = "6657"
)

func main() {
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Errorf("ERROR failed to listen:%v", err)
		return
	}
	s := grpc.NewServer()
	login_proto.RegisterLoginServer(s, &server{})
	err = s.Serve(lis)
	defer func() {
		s.Stop()
		lis.Close()
	}()
	if err != nil {
		log.Errorf("ERROR failed to start svr:%v", err)
	}
}

// GetMachineVerify 获取人机验证码
func (s *server) GetMachineVerify(ctx context.Context, req *login_proto.GetMachineVerifyReq) (
	*login_proto.GetMachineVerifyRsp, error) {
	rsp := &login_proto.GetMachineVerifyRsp{CodeId: "", Base64Img: "", RetCode: comm.GetVerCodeErr}
	capId, img, retCode, err := getMachineVerify(ctx)
	if err != nil || retCode != comm.SuccessCode {
		log.Errorf("get machine verify retCode:%v err:%v", retCode, err)
		return rsp, err
	}
	rsp.CodeId = capId
	rsp.Base64Img = img
	rsp.RetCode = comm.SuccessCode
	return rsp, nil
}

// 发送验证码结果
func (s *server) SendMachineVerifyResult(ctx context.Context, req *login_proto.SendMachineVerifyResultReq) (
	*login_proto.SendMachineVerifyResultRsp, error) {
	var rspErr error
	rsp := &login_proto.SendMachineVerifyResultRsp{RetCode: comm.VerifyCodeErr}
	retCode, err := sendMachineVerifyResult(ctx, req.CodeId, req.Ans)
	if err != nil {
		rspErr = fmt.Errorf("verify code err:%v", err)
		log.Errorf("%v", rspErr)
		return rsp, rspErr
	}
	rsp.RetCode = retCode
	return rsp, nil
}

// 发送一个短信验证码
func (s *server) SendTextVerCode(ctx context.Context, req *login_proto.SendTextVerCodeReq) (
	*login_proto.SendTextVerCodeRsp, error) {
	var rspErr error
	retCode, err := sendTextVerCode(ctx, req.PhoneNumber)
	if err != nil || retCode != comm.SuccessCode {
		rspErr = fmt.Errorf("send text vercode fail. retCode:%v err:%v", retCode, err)
		log.Errorf("%v", rspErr)
	}
	rsp := &login_proto.SendTextVerCodeRsp{RetCode: retCode}
	return rsp, rspErr
}

// 用户注册
func (s *server) UserRegister(ctx context.Context, req *login_proto.UserRegisterReq) (
	*login_proto.UserRegisterRsp, error) {
	var rspErr error
	retCode, err := register(ctx, req.PhoneNumber, req.Passwd, req.VerCode)
	if err != nil || retCode != comm.SuccessCode {
		rspErr = fmt.Errorf("register fail.retCode:%v err:%v req:%v", retCode, err, req)
		log.Errorf("%v", rspErr)
	}
	rsp := &login_proto.UserRegisterRsp{RetCode: retCode}
	return rsp, rspErr
}

// 用户密码登录
func (s *server) UserPswdLogin(ctx context.Context, req *login_proto.UserPswdLoginReq) (
	*login_proto.UserPswdLoginRsp, error) {
	var rspErr error
	retCode, err := passwordLogin(ctx, req.PhoneNumber, req.A1)
	if err != nil || retCode != comm.SuccessCode {
		rspErr = fmt.Errorf("password login fail. retCode:%v err:%v", retCode, err)
		log.Errorf("%v", rspErr)
	}
	rsp := &login_proto.UserPswdLoginRsp{RetCode: retCode}
	return rsp, rspErr
}

// 用户手机号码登录
func (s *server) UserPhoneLogin(ctx context.Context, req *login_proto.UserPhoneLoginReq) (
	*login_proto.UserPhoneLoginRsp, error) {
	var rspErr error
	retCode, err := verifyPhoneCode(ctx, req.PhoneNumber, req.VerCode)
	if err != nil || retCode != comm.SuccessCode {
		rspErr = fmt.Errorf("phone login fail. retCode:%v err:%v", retCode, err)
		log.Errorf("%v", rspErr)
	}
	rsp := &login_proto.UserPhoneLoginRsp{RetCode: retCode}
	return rsp, rspErr
}

// 重置密码
func (s *server) ResetPswd(ctx context.Context, req *login_proto.ResetPswdReq) (
	*login_proto.ResetPswdRsp, error) {
	var rspErr error
	retCode, err := resetPswd(ctx, req.PhoneNumber, req.VerCode, req.NewPw)
	if err != nil || retCode != comm.SuccessCode {
		rspErr = fmt.Errorf("resetPswd fail. phone:%v retCode:%v err:%v", req.PhoneNumber, retCode, err)
		log.Errorf("%v", rspErr)
	}
	rsp := &login_proto.ResetPswdRsp{RetCode: retCode}
	return rsp, rspErr
}
