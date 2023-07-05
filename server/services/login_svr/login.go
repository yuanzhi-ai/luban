package main

import (
	"context"
	"fmt"

	"github.com/yuanzhi-ai/luban/go_proto/verify_proto"
	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/log"
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
