// 验证码生成库
package comm

import (
	"math/rand"
	"time"

	"github.com/yuanzhi-ai/luban/server/log"
)

const (
	defaultWidth = 20
	defaultHight = 7
	verifyLen    = 6
	defaultDpi   = 72
	chars        = "ABCDEFGHIJKMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	charsLen     = len(chars)
)

// 图像验证码
// W图像宽度，H图像高度
type CaptchaGenerator struct {
	w, h, codeLen int
	fontSize      float64
	dpi           int
}

// GetNewCaptchaGenerator 获取一个验证码生成器
func GetNewCaptchaGenerator() *CaptchaGenerator {
	cg := CaptchaGenerator{w: defaultWidth, h: defaultHight, dpi: defaultDpi}
	cg.init()
	return &cg
}

// GetVerCode 获取一个验证码
// 返回一个验证码的id,和base64编码后的验证码图像
// 验证码的id = md5(验证码答案+秘钥)
// 故此解密只要再使用用户答案做相同计算与id比较即可
// 返回 验证码id， 验证码的url，和是否出错
func (cg *CaptchaGenerator) GetCaptcha() (string, string, error) {
	// 生成一个随机验证码
	captcha := cg.getRandStr()
	skeyInstance := GetSkeyInstance()
	capSkey, err := skeyInstance.GetSkey(CaptchaSkey)
	if err != nil || capSkey == "" {
		log.Errorf("get capSkey err:%v", err)
		return "", "", err
	}
	// 使用md5加密验证码和验证码的skey，作为改验证码的id
	capId := Md5Encode(captcha + capSkey)
	return capId, "", nil
	// 画图
}

// 初始化，设置随机数种子
func (cg *CaptchaGenerator) init() {
	rand.Seed(time.Now().UnixNano())
}

// getRandStr 获取一个验证码
func (cg *CaptchaGenerator) getRandStr() (randStr string) {
	for i := 0; i < cg.codeLen; i++ {
		randIndex := rand.Intn(charsLen)
		randStr += chars[randIndex : randIndex+1]
	}
	return randStr
}
