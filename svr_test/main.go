// Package http_entry 服务入口
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"google.golang.org/protobuf/proto"

	"github.com/yuanzhi-ai/luban/go_proto/login_proto"
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

func readJson(rsp interface{}, name string) error {
	path := fmt.Sprintf("./%v.json", name)
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("not find %v err:%v", path, err)
	}
	defer f.Close()
	jsonParser := json.NewDecoder(f)
	err = jsonParser.Decode(rsp)
	if err != nil {
		return fmt.Errorf("jsonParser.Decode %v fail. err:%v", path, err)
	}
	return nil
}

// DoRsp 回包
func DoRsp(w http.ResponseWriter, rsp proto.Message) {
	r, err := proto.Marshal(rsp)
	if err != nil {
		fmt.Printf("proto.Marshal err: %+v", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type, Token, token")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token, token")
	_, err = w.Write(r)
	if err != nil {
		fmt.Printf("w.Write err: %+v", err)
	}

}

// say hellow接口
func handler(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.EmptyReq{}
	DoRsp(w, rsp)
}

// getMachineVerifyHandler 获取人机验证码
func getMachineVerifyHandler(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.GetMachineVerifyRsp{}
	readJson(rsp, "GetMachineVerifyRsp")
	fmt.Printf("rsp:%+v", rsp)
	DoRsp(w, rsp)
}

// sendMachineVerifyResultHandler 人机验证码回包
func sendMachineVerifyResultHandler(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.SendMachineVerifyResultRsp{}
	readJson(rsp, "SendMachineVerifyResultRsp")
	DoRsp(w, rsp)
}

// 发送短信验证码
func sendTextVerCode(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.SendTextVerCodeRsp{}
	readJson(rsp, "SendTextVerCodeRsp")
	DoRsp(w, rsp)
}

// 用户账号密码登录
func userPswdLogin(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.UserPswdLoginRsp{}
	readJson(rsp, "UserPswdLoginRsp")
	DoRsp(w, rsp)
}

// 用户手机号登录
func userPhoneLogin(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.UserPhoneLoginRsp{}
	readJson(rsp, "UserPhoneLoginRsp")
	DoRsp(w, rsp)
}

// 用户注册 注册成功直接跳转登录
func userRegister(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.UserRegisterRsp{}
	readJson(rsp, "UserRegisterRsp")
	DoRsp(w, rsp)

}

// 用户重置密码
func resetPswd(w http.ResponseWriter, r *http.Request) {
	rsp := &login_proto.ResetPswdRsp{}
	readJson(rsp, "ResetPswdRsp")
	DoRsp(w, rsp)
}
