package svc

import (
	"grpc-common/ucenter/ucclient"
	"ucenter-api/internal/config"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	UCRegisterRpc ucclient.Register
	UCLoginRpc    ucclient.Login
	UCAssetRpc    ucclient.Asset
	UCMemberRpc   ucclient.Member
	UCWithdrawRpc ucclient.Withdraw
	// MarketRpc     mclient.Market
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		UCRegisterRpc: ucclient.NewRegister(zrpc.MustNewClient(c.UCenterRpc)),
		UCLoginRpc:    ucclient.NewLogin(zrpc.MustNewClient(c.UCenterRpc)),
		UCAssetRpc:    ucclient.NewAsset(zrpc.MustNewClient(c.UCenterRpc)),
		UCMemberRpc:   ucclient.NewMember(zrpc.MustNewClient(c.UCenterRpc)),
		UCWithdrawRpc: ucclient.NewWithdraw(zrpc.MustNewClient(c.UCenterRpc)),
		// MarketRpc:     mclient.NewMarket(zrpc.MustNewClient(c.MarketRpc)),
	}
}
