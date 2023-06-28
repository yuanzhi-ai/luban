// 公共逻辑类
package comm

import (
	"fmt"
	"sync"

	"github.com/coocood/freecache"
	"github.com/yuanzhi-ai/luban/server/log"
	"google.golang.org/protobuf/proto"
)

type MyCache struct {
	cache         *freecache.Cache
	expireSeconds int        //过期时间
	lock          sync.Mutex //锁
}

func (c *MyCache) Init(memSize int, timeOut int) {
	c.cache = freecache.NewCache(memSize)
	c.expireSeconds = int(timeOut)
}

// Get 获取一个数据
func (c *MyCache) Get(key []byte, value proto.Message) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	bufValue, err := c.cache.Get(key)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(bufValue, value)
	if err != nil {
		return fmt.Errorf("proto.Unmarshal err: %+v", err)
	}
	return nil
}

// 添加一个数据
func (c *MyCache) Set(key []byte, value proto.Message) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	bValue, err := proto.Marshal(value)
	if err != nil {
		log.Errorf("proto.Marshal err: %+v", err)
		return nil
	}
	err = c.cache.Set(key, bValue, c.expireSeconds)
	if err != nil {
		return err
	}
	return nil
}
