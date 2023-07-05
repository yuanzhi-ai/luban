// Package http_entry 服务入口
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yuanzhi-ai/luban/go_proto/login_proto"
	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/log"
	"github.com/yuanzhi-ai/luban/server/transer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/get-machine-verify", getMachineVerifyHandler)
	http.HandleFunc("/api/send-machine-verify-result", sendMachineVerifyResultHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("http.ListenAndServe err: ", err)
		return
	}
}

// say hellow接口
func handler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into say hello")
	rsp := &login_proto.EmptyReq{}
	transer.DoRsp(w, rsp)
	log.Infof("say hello ok!")

}

// getMachineVerifyHandler 获取人机验证码
func getMachineVerifyHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into getMachineVerifyHandler")
	req := &login_proto.GetMachineVerifyReq{}
	rsp := &login_proto.GetMachineVerifyRsp{RetCode: comm.GeneratorCaptchaErr}
	defer func() {
		transer.DoRsp(w, rsp)
	}()
	err := transer.GetReq(r, req)
	if err != nil {
		log.Errorf("trans GetMachineVerifyReq req err:%v", err)
		return
	}
	grpcConn, err := grpc.Dial("172.31.66.86:6657", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Dial login svr err:%v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client := login_proto.NewLoginClient(grpcConn)
	defer grpcConn.Close()
	defer cancel()
	rsp, err = client.GetMachineVerify(ctx, req)
	if err != nil {
		log.Errorf("RPC GetMachineVerify fial err:%v", err)
	}
	log.Debugf("req:%+v rsp:%+v", req, rsp)
}

func sendMachineVerifyResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into sendMachineVerifyResultHandler")
	req := &login_proto.SendMachineVerifyResultReq{}
	rsp := &login_proto.SendMachineVerifyResultRsp{RetCode: comm.VerifyCodeErr}
	defer func() {
		transer.DoRsp(w, rsp)
	}()
	err := transer.GetReq(r, req)
	if err != nil {
		log.Errorf("trans SendMachineVerifyResultReq req err:%v", err)
		return
	}
	grpcConn, err := grpc.Dial("172.31.66.86:6657", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Errorf("Dial login svr err:%v", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client := login_proto.NewLoginClient(grpcConn)
	defer grpcConn.Close()
	defer cancel()
	rsp, err = client.SendMachineVerifyResult(ctx, req)
	if err != nil {
		log.Errorf("RPC SendMachineVerifyResult err:%v", err)
		return
	}
	log.Debugf("req:%+v rsp:%+v", req, rsp)

}
