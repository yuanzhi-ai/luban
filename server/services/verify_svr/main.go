package main

import (
	"context"
	"net"

	"github.com/yuanzhi-ai/luban/go_proto/verify_proto"
	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/log"
	"google.golang.org/grpc"
)

type server struct {
	verify_proto.UnimplementedVerifyServer
}

const (
	PORT = "3360"
)

func main() {
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Errorf("ERROR failed to listen:%v", err)
		return
	}
	s := grpc.NewServer()
	verify_proto.RegisterVerifyServer(s, &server{})
	err = s.Serve(lis)
	defer func() {
		s.Stop()
		lis.Close()
	}()
	if err != nil {
		log.Errorf("ERROR failed to start svr:%v", err)
	}
}

// GetVerCode 获取一个验证码
func (s *server) GetVerCode(ctx context.Context, req *verify_proto.GetVerCodeReq) (
	*verify_proto.GetVerCodeRsp, error) {
	capId, encodeImg, err := getVerCode(ctx)
	if err != nil || capId == "" || encodeImg == "" {
		log.Errorf("getVerCode err:%v", err)
		rsp := &verify_proto.GetVerCodeRsp{CodeId: "", Base64Img: "", RetCode: comm.GeneratorCaptchaErr}
		return rsp, err
	}
	rsp := &verify_proto.GetVerCodeRsp{CodeId: capId, Base64Img: encodeImg, RetCode: 0}
	return rsp, nil
}

// VerifyCode 验证用户的验证码
func (s *server) VerifyCode(ctx context.Context, req *verify_proto.VerifyCodeReq) (
	*verify_proto.VerifyCodeRsp, error) {
	// 默认验证失败
	rsp := &verify_proto.VerifyCodeRsp{RetCode: comm.VerifyCaptchaErr}
	success, retCode, err := verifyCode(ctx, req.CodeId, req.Ans)
	if err != nil || retCode != comm.SuccessCode || !success {
		log.Errorf("verify code err:%v", err)
		return rsp, err
	}
	rsp.RetCode = comm.SuccessCode
	return rsp, nil
}
