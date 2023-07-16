package main

import (
	"context"
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
	rsp := &login_proto.SendMachineVerifyResultRsp{RetCode: comm.VerifyCodeErr}
	retCode, err := sendMachineVerifyResult(ctx, req.CodeId, req.Ans)
	if err != nil {
		log.Errorf("verify code err:%v", err)
		return rsp, err
	}
	rsp.RetCode = retCode
	return rsp, nil
}

// 发送一个短信验证码
func (s *server) SendTextVerCode(ctx context.Context, req *login_proto.SendTextVerCodeReq) (
	*login_proto.SendTextVerCodeRsp, error) {
	return nil, nil
}

func (s *server) UserRegister(ctx context.Context, req *login_proto.UserRegisterReq) (
	*login_proto.UserRegisterRsp, error) {
	return nil, nil
}

func (s *server) UserPswdLogin(ctx context.Context, req *login_proto.UserPswdLoginReq) (
	*login_proto.UserPswdLoginRsp, error) {
	return nil, nil
}

func (s *server) UserPhoneLogin(ctx context.Context, req *login_proto.UserPhoneLoginReq) (
	*login_proto.UserPhoneLoginRsp, error) {
	return nil, nil
}

func (s *server) ResetPswd(ctx context.Context, req *login_proto.ResetPswdReq) (
	*login_proto.ResetPswdRsp, error) {
	return nil, nil
}
