package main

import (
	"context"
	"fmt"

	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/log"
)

// getVerCode 获取验证码真正的执行
func getVerCode(ctx context.Context) (string, string, error) {
	captchaGenerator := comm.GetNewCaptchaGenerator()
	capId, encodeImg, err := captchaGenerator.GetCaptcha()
	if capId == "" || encodeImg == "" || err != nil {
		log.Errorf("get capcha err:%v", err)
		return "", "", err
	}
	return capId, encodeImg, nil
}

// verifyCode 验证用户答案
func verifyCode(ctx context.Context, capId string, userAnswer string) (bool, int32, error) {
	if len(capId) < 1 || len(capId) > 100 || len(userAnswer) != comm.VerifyLen {
		return false, comm.InputErr, fmt.Errorf(
			"error input params. len(capId):%v len(userAnswer):%v", len(capId), len(userAnswer))
	}
	captchaGenerator := comm.GetNewCaptchaGenerator()
	success, err := captchaGenerator.VerifyCode(capId, userAnswer)
	if err != nil {
		return false, comm.VerifyCaptchaErr, err
	}
	return success, comm.SuccessCode, nil
}
