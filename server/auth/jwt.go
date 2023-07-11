// 生成jwt 签名
package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yuanzhi-ai/luban/server/comm"
	"github.com/yuanzhi-ai/luban/server/log"
)

const (
	MachineJwtType = "machineAuth"
)

// Jwt 类型的接口
type Jwt interface {
	generatorJWT(payload map[string]interface{}) (string, error)
}

// JwtHead jwt的头
var JwtHead string

func initJwtHead() error {
	if JwtHead == "" {
		jwtHead := map[string]string{"alg": "md5", "type": "JWT"}
		jsonHead, err := json.Marshal(jwtHead)
		if err != nil {
			log.Errorf("json marshal jwt head:%v err:%v", jwtHead, err)
			return err
		}
		JwtHead = base64.StdEncoding.EncodeToString(jsonHead)
	}
	return nil
}

// jwtEncode 对消息体进行jwt编码
func jwtEncode(payload map[string]interface{}) (string, error) {
	err := initJwtHead()
	if err != nil {
		log.Errorf("init jwt head err:%v", err)
		return "", err
	}
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

// 通用jwt结构
type GeneralJwt struct {
}

// GeneratorJWT 通用的Jwt签名
func (j *GeneralJwt) generatorJWT(payload map[string]interface{}) (string, error) {
	jwtSignature, err := jwtEncode(payload)
	if err != nil {
		log.Errorf("jwt encode err:%v", err)
		return "", err
	}
	return jwtSignature, nil
}

// 人机验证的jwt结构
type MachineJwt struct {
}

func (j *MachineJwt) generatorJWT(payload map[string]interface{}) (string, error) {
	if _, ok := payload["owner"]; !ok {
		return "", fmt.Errorf("chekparam err: no owner")
	}
	nowTime := time.Now().Second()
	expTime := nowTime + 5*60
	defPayload := map[string]interface{}{"iat": nowTime, "exp": expTime, "purpose": MachineJwtType, "owner": payload["owner"]}
	jwtSingature, err := jwtEncode(defPayload)
	if err != nil {
		log.Errorf("jwt encode err:%v", err)
		return "", err
	}
	return jwtSingature, nil
}

func GeneratorJwt(jwtType string, payload map[string]interface{}) (string, error) {
	var jwt Jwt
	switch jwtType {
	case MachineJwtType:
		jwt = &MachineJwt{}
	default:
		jwt = &GeneralJwt{}
	}
	return jwt.generatorJWT(payload)
}
