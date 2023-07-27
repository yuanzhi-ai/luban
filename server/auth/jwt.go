// 生成jwt 签名
package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/log"
)

// Jwt 类型的接口
type Jwt interface {
	generatorJWT(payload map[string]interface{}) (string, error)
}

// JwtHead jwt的头
var JwtHead string

func init() {
	jwtHeadTmp := map[string]string{"alg": "md5", "type": "JWT"}
	jsonHead, _ := json.Marshal(jwtHeadTmp)
	JwtHead = base64.StdEncoding.EncodeToString(jsonHead)
}

// jwtEncode 对消息体进行jwt编码
func jwtEncode(payload map[string]interface{}) (string, error) {
	// 对消息体编码
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Errorf("json marshal payload:%+v err:%v", jsonPayload, err)
		return "", err
	}
	jwtPayload := base64.StdEncoding.EncodeToString(jsonPayload)
	// 对消息头和消息体签名
	value := fmt.Sprintf("%v.%v", JwtHead, jwtPayload)
	skeyInstance := comm.GetSkeyInstance()
	jwtSkey, err := skeyInstance.GetSkey(comm.JwtSkey)
	if err != nil || jwtSkey == "" {
		log.Errorf("get jwtSkey err%v", err)
		return "", err
	}
	// 这里对value拼接skey做md5加密
	signature := comm.Md5Encode(value + jwtSkey)
	jwtSignature := fmt.Sprintf("%v.%v", value, signature)
	return jwtSignature, nil
}

func GeneratorJWT(owner string, api string, expTs int64) (string, error) {
	nowTime := time.Now().Unix()
	expTime := fmt.Sprintf("%v", (nowTime + expTs))
	payload := map[string]interface{}{"iat": fmt.Sprintf("%v", nowTime), "exp": expTime, "purpose": api, "owner": owner}
	jwtSingature, err := jwtEncode(payload)
	if err != nil {
		log.Errorf("jwt encode err:%v", err)
		return "", err
	}
	return jwtSingature, nil

}

// 登录态续期
func JwtExtension(payload map[string]interface{}, api string, expTs int64) (string, error) {
	nowTime := time.Now().Unix()
	expTime := nowTime + expTs
	payload["iat"] = nowTime
	payload["exp"] = expTime
	jwtSingature, err := jwtEncode(payload)
	if err != nil {
		log.Errorf("jwt encode err:%v", err)
		return "", err
	}
	return jwtSingature, nil
}

// 判断jwt签名是否合法
func IsJwtSignatureLegal(jwt string) bool {
	values := strings.Split(jwt, ".")
	// jwt3段式
	if len(values) != 3 {
		return false
	}
	head := values[0]
	payload := values[1]
	signature := values[2]
	// 检查签名头
	if head != JwtHead {
		return false
	}
	// 验证签名是否一致
	skeyInstance := comm.GetSkeyInstance()
	jwtSkey, err := skeyInstance.GetSkey(comm.JwtSkey)
	if err != nil || jwtSkey == "" {
		log.Errorf("get jwtSkey err%v", err)
		return false
	}
	svrSignature := comm.Md5Encode(fmt.Sprintf("%v.%v", head, payload) + jwtSkey)
	if svrSignature != signature {
		log.Errorf("forged jwt! jwt:%v", jwt)
		return false
	}
	return true
}

// 解码jwt
func JwtDecodePayload(jwt string) (map[string]interface{}, error) {
	if !IsJwtSignatureLegal(jwt) {
		return nil, fmt.Errorf("jwt err")
	}
	d, err := base64.StdEncoding.DecodeString(strings.Split(jwt, ".")[1])
	if err != nil {
		return nil, fmt.Errorf("jwt payload decode err:%v", err)
	}
	var payload map[string]interface{}
	err = json.Unmarshal(d, &payload)
	if err != nil {
		return nil, fmt.Errorf("json unmarshal jwt paylod fail. json:%v err:%v", string(d), err)
	}
	return payload, nil
}
