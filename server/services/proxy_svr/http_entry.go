// Package http_entry 服务入口
package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/yuanzhi-ai/luban/go_proto/login_proto"
	"github.com/yuanzhi-ai/luban/server/auth"
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
	http.HandleFunc("/api/send-text-ver-code", sendTextVerCode)
	http.HandleFunc("/api/user-pswd-login", userPswdLogin)
	http.HandleFunc("/api/user-phone-login", userPhoneLogin)
	http.HandleFunc("/api/user-register", userRegister)
	http.HandleFunc("/api/reset-pswd", resetPswd)

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
		log.Debugf("req:%+v rsp:%+v", req, rsp)
		transer.DoRsp(w, rsp)
	}()
	if r.Method == "OPTIONS" {
		rsp.RetCode = comm.SuccessCode
		return
	}
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
	if err != nil || rsp.RetCode != comm.SuccessCode {
		log.Errorf("RPC GetMachineVerify fial err:%v", err)
	}
	rsp.RetCode = comm.SuccessCode

}

// sendMachineVerifyResultHandler 人机验证码回包
func sendMachineVerifyResultHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into sendMachineVerifyResultHandler")
	req := &login_proto.SendMachineVerifyResultReq{}
	rsp := &login_proto.SendMachineVerifyResultRsp{RetCode: comm.VerifyCodeErr}
	defer func() {
		log.Debugf("req:%+v rsp:%+v", req, rsp)
		transer.DoRsp(w, rsp)
	}()
	if r.Method == "OPTIONS" {
		rsp.RetCode = comm.SuccessCode
		return
	}

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
	if err != nil || rsp.RetCode != comm.SuccessCode {
		log.Errorf("RPC SendMachineVerifyResult err:%v", err)
		return
	}

	// 验证成功的带上jwt做游客态签名
	jwt, err := auth.GeneratorJWT(req.CodeId, auth.PRCGetMachineVerify, auth.TouristJwtExpTs)
	if err != nil {
		log.Errorf("GeneratorJwt err:%v", err)
		return
	}
	rsp.RetCode = comm.SuccessCode
	w.Header().Set("Token", jwt)
}

// 发送短信验证码
func sendTextVerCode(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into sendTextVerCode")
	req := &login_proto.SendTextVerCodeReq{}
	rsp := &login_proto.SendTextVerCodeRsp{RetCode: comm.SendSmsVerCodeErr}
	defer func() {
		log.Debugf("req:%+v rsp:%+v", req, rsp)
		transer.DoRsp(w, rsp)
	}()
	if r.Method == "OPTIONS" {
		rsp.RetCode = comm.SuccessCode
		return
	}
	// 先做jwt的校验
	jwt := r.Header.Get("Token")
	log.Debugf("r.Header.Get Token:%v", jwt)
	payload, err := auth.JwtDecodePayload(jwt)
	if err != nil || payload == nil || len(payload) == 0 {
		log.Errorf("jwt decode payload jwt:%v err:%v", jwt, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查payload是否合法
	err = IsPayloadLegal(payload, auth.PRCGetMachineVerify)
	if err != nil {
		log.Errorf("jwt payload legal. payload:%+v err:%v", payload, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查用户访问限频
	if !auth.FrequencyControler.CanVisit(req.PhoneNumber, auth.RPCSendTextVerCode) {
		rsp.RetCode = comm.FrequencyErr
		return
	}
	// rpc请求
	err = transer.GetReq(r, req)
	if err != nil {
		log.Errorf("trans SendTextVerCodeReq req:%v err:%v", req, err)
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
	rsp, err = client.SendTextVerCode(ctx, req)
	if err != nil || rsp == nil || rsp.RetCode != comm.SuccessCode {
		log.Errorf("RPC SendTextVerCode err:%v", err)
		return
	}
	// 成功后下发短信验证成功的jwt
	jwt, err = auth.GeneratorJWT(req.PhoneNumber, auth.RPCSendTextVerCode, auth.TouristJwtExpTs)
	if err != nil || jwt == "" {
		log.Errorf("generator jwt err:%v", err)
		return
	}
	w.Header().Set("Token", jwt)
	rsp.RetCode = comm.SuccessCode
}

// 用户账号密码登录
func userPswdLogin(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into userPswdLogin")
	req := &login_proto.UserPswdLoginReq{}
	rsp := &login_proto.UserPswdLoginRsp{RetCode: comm.LoginErr}
	defer func() {
		log.Debugf("req:%+v rsp:%+v", req, rsp)
		transer.DoRsp(w, rsp)
	}()
	if r.Method == "OPTIONS" {
		rsp.RetCode = comm.SuccessCode
		return
	}
	// 先做jwt的校验
	jwt := r.Header.Get("Token")
	payload, err := auth.JwtDecodePayload(jwt)
	if err != nil || payload == nil || len(payload) == 0 {
		log.Errorf("jwt decode payload jwt:%v err:%v", jwt, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查payload jwt是否正确，验证码是否过期
	err = IsPayloadLegal(payload, auth.RPCSendMachineVerifyResult)
	if err != nil {
		log.Errorf("jwt payload legal. payload:%+v err:%v", payload, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查限频
	if !auth.FrequencyControler.CanVisit(req.PhoneNumber, auth.RPCUserPswdLogin) {
		rsp.RetCode = comm.FrequencyErr
		return
	}
	//rpc访问
	err = transer.GetReq(r, req)
	if err != nil {
		log.Errorf("trans UserPhoneLoginReq req:%v err:%v", req, err)
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
	rsp, err = client.UserPswdLogin(ctx, req)
	if err != nil || rsp.RetCode != comm.SuccessCode {
		log.Errorf("RPC UserPswdLogin err:%v", err)
		return
	}
	// 成功后下发登录的jwt
	jwt, err = auth.GeneratorJWT(req.PhoneNumber, auth.JwtLoginPurpose, auth.HomePageJwtExpTs)
	if err != nil || jwt == "" {
		log.Errorf("generator jwt err:%v", err)
		return
	}
	w.Header().Set("Token", jwt)
	rsp.RetCode = comm.SuccessCode
}

// 用户手机号登录
func userPhoneLogin(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into userPhoneLogin")
	req := &login_proto.UserPhoneLoginReq{}
	rsp := &login_proto.UserPhoneLoginRsp{RetCode: comm.LoginErr}
	defer func() {
		log.Debugf("req:%+v rsp:%+v", req, rsp)
		transer.DoRsp(w, rsp)
	}()
	if r.Method == "OPTIONS" {
		rsp.RetCode = comm.SuccessCode
		return
	}
	// 先做jwt的校验
	jwt := r.Header.Get("Token")
	payload, err := auth.JwtDecodePayload(jwt)
	if err != nil || payload == nil || len(payload) == 0 {
		log.Errorf("jwt decode payload jwt:%v err:%v", jwt, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查payload jwt是否正确，验证码是否过期
	err = IsPayloadLegal(payload, auth.RPCSendTextVerCode)
	if err != nil {
		log.Errorf("jwt payload legal. payload:%+v err:%v", payload, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查限频
	if !auth.FrequencyControler.CanVisit(req.PhoneNumber, auth.RPCUserPhoneLogin) {
		rsp.RetCode = comm.FrequencyErr
		return
	}
	//rpc访问
	err = transer.GetReq(r, req)
	if err != nil {
		log.Errorf("trans UserPhoneLoginReq req:%v err:%v", req, err)
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
	rsp, err = client.UserPhoneLogin(ctx, req)
	if err != nil || rsp.RetCode != comm.SuccessCode {
		log.Errorf("RPC UserPswdLogin err:%v", err)
		return
	}
	// 成功后下发登录的jwt
	jwt, err = auth.GeneratorJWT(req.PhoneNumber, auth.JwtLoginPurpose, auth.HomePageJwtExpTs)
	if err != nil || jwt == "" {
		log.Errorf("generator jwt err:%v", err)
		return
	}
	w.Header().Set("Token", jwt)
	rsp.RetCode = comm.SuccessCode
}

// 用户注册 注册成功直接跳转登录
func userRegister(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into userRegister")
	req := &login_proto.UserRegisterReq{}
	rsp := &login_proto.UserRegisterRsp{RetCode: comm.RegisterErr}
	defer func() {
		log.Debugf("req:%+v rsp:%+v", req, rsp)
		transer.DoRsp(w, rsp)
	}()
	if r.Method == "OPTIONS" {
		rsp.RetCode = comm.SuccessCode
		return
	}
	// 先做jwt的校验
	jwt := r.Header.Get("Token")
	payload, err := auth.JwtDecodePayload(jwt)
	if err != nil || payload == nil || len(payload) == 0 {
		log.Errorf("jwt decode payload jwt:%v err:%v", jwt, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查payload jwt是否正确，验证码是否过期
	err = IsPayloadLegal(payload, auth.RPCSendTextVerCode)
	if err != nil {
		log.Errorf("jwt payload legal. payload:%+v err:%v", payload, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查限频
	if !auth.FrequencyControler.CanVisit(req.PhoneNumber, auth.RPCUserRegister) {
		rsp.RetCode = comm.FrequencyErr
		return
	}
	//rpc访问
	err = transer.GetReq(r, req)
	if err != nil {
		log.Errorf("trans UserRegisterReq req:%v err:%v", req, err)
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
	rsp, err = client.UserRegister(ctx, req)
	if err != nil || rsp.RetCode != comm.SuccessCode {
		log.Errorf("RPC UserPswdLogin err:%v", err)
		return
	}
	// 成功后下发登录的jwt
	jwt, err = auth.GeneratorJWT(req.PhoneNumber, auth.JwtLoginPurpose, auth.HomePageJwtExpTs)
	if err != nil || jwt == "" {
		log.Errorf("generator jwt err:%v", err)
		return
	}
	w.Header().Set("Token", jwt)
	rsp.RetCode = comm.SuccessCode
}

// 用户重置密码
func resetPswd(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into resetPswd")
	req := &login_proto.ResetPswdReq{}
	rsp := &login_proto.ResetPswdRsp{RetCode: comm.LoginErr}
	defer func() {
		log.Debugf("req:%+v rsp:%+v", req, rsp)
		transer.DoRsp(w, rsp)
	}()
	if r.Method == "OPTIONS" {
		rsp.RetCode = comm.SuccessCode
		return
	}
	// 先做jwt的校验
	jwt := r.Header.Get("Token")
	payload, err := auth.JwtDecodePayload(jwt)
	if err != nil || payload == nil || len(payload) == 0 {
		log.Errorf("jwt decode payload jwt:%v err:%v", jwt, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查payload jwt是否正确，验证码是否过期
	err = IsPayloadLegal(payload, auth.RPCSendTextVerCode)
	if err != nil {
		log.Errorf("jwt payload legal. payload:%+v err:%v", payload, err)
		rsp.RetCode = comm.JWTErr
		return
	}
	// 检查限频
	if !auth.FrequencyControler.CanVisit(req.PhoneNumber, auth.RPCResetPswd) {
		rsp.RetCode = comm.FrequencyErr
		return
	}
	//rpc访问
	err = transer.GetReq(r, req)
	if err != nil {
		log.Errorf("trans ResetPswdReq req:%v err:%v", req, err)
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
	rsp, err = client.ResetPswd(ctx, req)
	if err != nil || rsp.RetCode != comm.SuccessCode {
		log.Errorf("RPC ResetPswd err:%v", err)
		return
	}
	// 这里成功后跳主页，重走人机验证，不下发jwt
	rsp.RetCode = comm.SuccessCode
}
