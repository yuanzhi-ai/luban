package data

import (
	"context"
	"fmt"
	"time"

	"github.com/yuanzhi-ai/luban/server/log"
	"github.com/yuanzhi-ai/luban/server/repo/mredis"
)

const (
	// 电话验证码过期时间
	PhoneCodeExpTime = time.Duration(300) * time.Second
	NotFoundKey      = "not found key"
)

func phoneCodeKey(phone string) string {
	return fmt.Sprintf("phone_code_%v", phone)
}

// SetPhoneCaptcha 电话验证码写入redis
func SetPhoneCaptcha(ctx context.Context, phone string, code string) error {
	rc := mredis.GetRdbClient()
	rkey := phoneCodeKey(phone)
	err := rc.Set(ctx, rkey, code, PhoneCodeExpTime).Err()
	return err
}

// GetPhoneCaptcha 获取手机号的验证码
func GetPhoneCaptcha(ctx context.Context, phone string) (string, error) {
	rc := mredis.GetRdbClient()
	rkey := phoneCodeKey(phone)
	val, err := rc.Get(ctx, rkey).Result()
	if err != nil && err != mredis.EmptyKeyErr {
		return "", err
	}
	if err == mredis.EmptyKeyErr {
		return "", err
	}
	return val, nil
}

// DelPhoneCaptcha 删除手机验证码
func DelPhoneCaptcha(ctx context.Context, phone string) {
	rc := mredis.GetRdbClient()
	rkey := phoneCodeKey(phone)
	val, err := rc.Do(ctx, "del", rkey).Result()
	if err != nil {
		log.Errorf("del key:%v fial. err:%v", rkey, err)
	}
	log.Debugf("del key:%v return value:%v", rkey, val)
}
