package main

import (
	"fmt"
	"strconv"
	"time"
)

// checkTextVerJwt 检测发送短信验证码的jwt登录态
func IsPayloadLegal(payload map[string]interface{}, api string) error {
	// 检查jwt是否过期
	exp, ok := payload["exp"]
	if !ok {
		return fmt.Errorf("key exp not found")
	}
	expTs, err := strconv.ParseInt(exp.(string), 10, 64)
	if err != nil {
		return fmt.Errorf("conv exp to ts fail. exp:%v", exp)
	}
	now := time.Now().Unix()
	if expTs > now {
		return fmt.Errorf("jwt expiration now:%v exp:%v", now, expTs)
	}
	// 检查purpose 是否正确
	purpos, ok := payload["purpose"]
	if !ok {
		return fmt.Errorf("key purpose not found")
	}
	purpos = purpos.(string)
	if purpos != api {
		return fmt.Errorf("purpost not equal api. purpos:%v api:%v", purpos, api)
	}
	// 检查是否有owner字段
	owmer, ok := payload["owner"]
	if !ok || owmer.(string) == "" {
		return fmt.Errorf("key owner not found. owner:%v", owmer)
	}
	return nil
}
