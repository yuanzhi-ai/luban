// 读取服务器各类秘钥
// 单例类
package comm

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/yuanzhi-ai/luban/server/log"

	"github.com/robfig/cron/v3"
)

const (
	// errSkey 读取skey失败时的标识
	errSkey = ""
	// CaptchaSkey 验证码秘钥存枚举
	CaptchaSkey     = "captcha"
	captchaSkeyPath = "./skey/captchaSkey.txt"
	// 登录注册的秘钥枚举
	LoginRegisterSkey     = "login"
	LoginRegisterSkeyPath = "./skey/loginSkey.txt"
)

type skey struct {
	// skeyFilePaths 秘钥保存的文件路径
	skeyFilePaths map[string]string
	// allSkeys 存放所有的skey
	allSkeys map[string]string
}

var instance *skey
var instanceLock sync.Mutex
var skeyLock sync.Mutex // skey的锁

// GetSkeyInstance 获取skey实例
func GetSkeyInstance() *skey {
	if instance == nil {
		instanceLock.Lock()
		instance = &skey{}
		instance.init()
		instanceLock.Unlock()
	}
	return instance
}

// 获取一个skey
func (s *skey) GetSkey(skeyMod string) (string, error) {
	if _, ok := s.allSkeys[skeyMod]; !ok {
		err := fmt.Errorf("skeyMod not found in allSkeys skeyMod:%v", skeyMod)
		log.Errorf("err:%v", err)
		return errSkey, err
	}
	return s.allSkeys[skeyMod], nil
}

func (s *skey) init() {
	// 这里初始化所有的秘钥到文件里
	s.skeyFilePaths = map[string]string{CaptchaSkey: captchaSkeyPath, LoginRegisterSkey: LoginRegisterSkeyPath}
	s.reloadAllSkey()
	// 每5分钟重新load一次秘钥
	c := cron.New()
	_, err := c.AddFunc("*/5 * * * *", reloadSkey)
	if err != nil {
		log.Errorf("启动skey定时刷新失败! err:%v", err)
	}
	c.Start()

}

func reloadSkey() {
	skeyLock.Lock()
	defer skeyLock.Unlock()
	instance := GetSkeyInstance()
	instance.reloadAllSkey()
}

// reloadAllSkey 重新加载所有的skey
func (s *skey) reloadAllSkey() {
	s.allSkeys = make(map[string]string)
	for skeyMod := range s.skeyFilePaths {
		s.readSkeyFile(skeyMod)
	}
}

// 从文件中读取一个skey
func (s *skey) readSkeyFile(skeyMod string) {

	fileSkey := errSkey
	defer func() {
		s.allSkeys[skeyMod] = fileSkey
	}()
	// skye的类型检查
	if _, ok := s.skeyFilePaths[skeyMod]; !ok {
		log.Errorf("fail open skey file, skey type not found. skey type:%v", skeyMod)
		return
	}
	// 文件读取
	path := s.skeyFilePaths[skeyMod]
	f, err := os.Open(path)
	if err != nil {
		log.Errorf("fail open skey file, file path:%v", path)
		return
	}
	defer f.Close()
	// 读取秘钥
	r := bufio.NewReader(f)
	bytes, _, err := r.ReadLine()
	if err != nil {
		log.Errorf("fail read skey file line, file path:%v err:%v", path, err)
	}
	fileSkey = string(bytes)
}
