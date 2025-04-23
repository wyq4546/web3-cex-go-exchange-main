package main

import (
	"flag"
	"grpc-common/ucenter/types/asset"
	"grpc-common/ucenter/types/login"
	"grpc-common/ucenter/types/member"
	"grpc-common/ucenter/types/register"
	"grpc-common/ucenter/types/withdraw"
	"ucenter/internal/config"
	"ucenter/internal/server"
	"ucenter/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()

	// 配置日志
	logx.MustSetup(logx.LogConf{
		ServiceName: "ucenter-rpc",
		Mode:        "console",
		Encoding:    "json",
		TimeFormat:  "2006-01-02 15:04:05",
		Level:       "info",
		Path:        "logs",
		Compress:    false,
		KeepDays:    7,
		Stat:        false,
	})

	logx.Info("Starting ucenter RPC server...")

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	logx.Info("Service context initialized")

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		logx.Info("Registering gRPC services...")

		register.RegisterRegisterServer(grpcServer, server.NewRegisterServer(ctx))
		logx.Info("Register service registered")

		login.RegisterLoginServer(grpcServer, server.NewLoginServer(ctx))
		logx.Info("Login service registered")

		asset.RegisterAssetServer(grpcServer, server.NewAssetServer(ctx))
		logx.Info("Asset service registered")

		member.RegisterMemberServer(grpcServer, server.NewMemberServer(ctx))
		logx.Info("Member service registered")

		withdraw.RegisterWithdrawServer(grpcServer, server.NewWithdrawServer(ctx))
		logx.Info("Withdraw service registered")


		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
			logx.Info("gRPC reflection registered (dev/test mode)")
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...", c.ListenOn)
	s.Start()
}
