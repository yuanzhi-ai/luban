// Package transer 解析及打包请求和回包
package transer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/yuanzhi-ai/luban/server/log"
	"google.golang.org/protobuf/proto"
)

// GetReq 从body中获取请求结构,req为req结构指针
func GetReq(r *http.Request, req proto.Message) error {

	body := bytes.NewBuffer(make([]byte, 0))
	_, err := io.Copy(body, r.Body)
	if err != nil {
		log.Errorf("read body err: %+v", err)
		return fmt.Errorf("read body err: %+v", err)
	}
	log.Debugf("req的bytes数组 body.Bytes():%+v", body.Bytes())
	err = proto.Unmarshal(body.Bytes(), req)
	if err != nil {
		log.Errorf("proto.Unmarshal err: %+v", err)
		return fmt.Errorf("proto.Unmarshal err: %+v", err)
	}
	return nil
}

// DoRsp 回包
func DoRsp(w http.ResponseWriter, rsp proto.Message) {
	r, err := proto.Marshal(rsp)
	if err != nil {
		log.Errorf("proto.Marshal err: %+v", err)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "content-type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
	_, err = w.Write(r)
	if err != nil {
		log.Errorf("w.Write err: %+v", err)
	}

}

// SetRetCode 设置返回码
func SetRetCode(w http.ResponseWriter, ret int) {
	if ret != http.StatusOK {
		ret = http.StatusExpectationFailed
	}
	w.WriteHeader(ret)
}
