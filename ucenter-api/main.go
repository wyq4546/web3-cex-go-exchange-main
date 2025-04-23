package main

import (
	"flag"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"

	"ucenter-api/internal/config"
	"ucenter-api/internal/svc"
	"ucenter-api/internal/handler"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/conf.yaml", "the config file")

func main() {
	flag.Parse()
	// 配置日志
	logx.MustSetup(logx.LogConf{
		ServiceName: "ucenter-rpc",
		Mode:        "console",
		Encoding:    "plain",
		TimeFormat:  "2006-01-02 15:04:05",
		Level:       "info",
		Path:        "logs",
		Compress:    false,
		KeepDays:    7,
		Stat:        false,
	})

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(func(header http.Header) {
		header.Set("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization,token,x-auth-token")
	}, nil, "http://localhost:8888"))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	router := handler.NewRouters(server)
	handler.RegisterHandlers(router, ctx)

	logx.Infof("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
