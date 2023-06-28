// Package http_entry 服务入口
package main

import (
	"fmt"
	"net/http"

	"github.com/yuanzhi-ai/luban/go_proto/login_proto"
	"github.com/yuanzhi-ai/luban/server/log"
	"github.com/yuanzhi-ai/luban/server/transer"
)

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("http.ListenAndServe err: ", err)
		return
	}
}

// say hellow接口
func handler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("into say hello")
	rsp := &login_proto.EmptyReq{}
	transer.DoRsp(w, rsp)
	log.Infof("say hello ok!")

}
