package main

import (
	"context"
	"net"

	"github.com/yuanzhi-ai/luban/go_proto/verify_proto"
	"github.com/yuanzhi-ai/luban/server/log"
	"google.golang.org/grpc"
)

type server struct {
	verify_proto.UnimplementedVerifyServer
}

const (
	PORT = "3360"
)

func main() {
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Errorf("ERROR failed to listen:%v", err)
		return
	}
	s := grpc.NewServer()
	verify_proto.RegisterVerifyServer(s, &server{})
	err = s.Serve(lis)
	defer func() {
		s.Stop()
		lis.Close()
	}()
	if err != nil {
		log.Errorf("ERROR failed to start svr:%v", err)
	}
}

// GetVerCode 获取一个验证码
func (s *server) GetVerCode(ctx context.Context, req *verify_proto.GetVerCodeReq) (
	*verify_proto.GetVerCodeRsp, error) {

	return nil, nil
}
