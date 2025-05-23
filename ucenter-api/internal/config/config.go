package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	UCenterRpc    zrpc.RpcClientConf
	MarketRpc     zrpc.RpcClientConf
	JWT           AuthConfig
	PlaygroundRpc zrpc.RpcClientConf
}
type AuthConfig struct {
	AccessSecret string
	AccessExpire int64
}
