/*接口限频*/
package auth

import (
	"sync"
	"time"
)

// FrequecnyControl 频控
type frequencyControl struct {
	counter map[string]int64 // 每个uin在一个周期内请求了多少次
	tsGap   int64            //请求间隔 即：每两次请求之间至少间隔多少秒
	locker  sync.Mutex       //锁
}

// canQuery 是否可以访问
func (f *frequencyControl) canQuery(id string) bool {
	f.locker.Lock()
	defer f.locker.Unlock()
	lastTs, ok := f.counter[id]
	now := time.Now().Unix()
	if !ok || now-lastTs > f.tsGap {
		f.counter[id] = now
		return true
	}
	return false
}

// GlobalFrequency 全局限速锁
type GlobalFrequency struct {
	control map[string]*frequencyControl
}

var FrequencyControler *GlobalFrequency

func init() {
	FrequencyControler = &GlobalFrequency{control: make(map[string]*frequencyControl)}
	FrequencyControler.control[RPCGetVerCode] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fGetVerCode, locker: sync.Mutex{}}
	FrequencyControler.control[RPCVerifyCode] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fVerifyCode, locker: sync.Mutex{}}
	FrequencyControler.control[PRCGetMachineVerify] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fGetMachineVerify, locker: sync.Mutex{}}
	FrequencyControler.control[RPCSendMachineVerifyResult] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fSendMachineVerifyResult, locker: sync.Mutex{}}
	FrequencyControler.control[RPCSendTextVerCode] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fSendTextVerCode, locker: sync.Mutex{}}
	FrequencyControler.control[RPCUserPswdLogin] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fUserPswdLogin, locker: sync.Mutex{}}
	FrequencyControler.control[RPCUserPhoneLogin] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fUserPhoneLogin, locker: sync.Mutex{}}
	FrequencyControler.control[RPCUserRegister] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fUserRegister, locker: sync.Mutex{}}
	FrequencyControler.control[RPCResetPswd] = &frequencyControl{counter: make(map[string]int64, defaultSize), tsGap: fResetPswd, locker: sync.Mutex{}}
}

func (g *GlobalFrequency) CanVisit(id string, api string) bool {
	f, ok := g.control[api]
	if !ok {
		return true
	}
	return f.canQuery(id)
}
