package main

import "sync"

type captchaCache struct {
	captchaSet map[string]bool
	locker     sync.RWMutex
	maxLen     int
}

var captchaSet captchaCache

func (c *captchaCache) init() {
	const cacheLimit = 1000000
	captchaSet = captchaCache{
		captchaSet: make(map[string]bool, cacheLimit),
		locker:     sync.RWMutex{},
		maxLen:     cacheLimit,
	}
}

// isCaptchaUsed 验证码是否已经用完了
func (c *captchaCache) isCaptchaUsed(capID string) bool {
	_, ok := captchaSet.captchaSet[capID]
	return ok
}

// addCaptcha 添加一个验证码
func (c *captchaCache) addCaptcha(capID string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	if len(c.captchaSet) >= c.maxLen {
		c.captchaSet = make(map[string]bool, c.maxLen)
	}
	c.captchaSet[capID] = true
}
