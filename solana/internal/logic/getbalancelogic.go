package logic

import (
	"context"

	"solana/internal/svc"
	"grpc-common/solana/types/playground"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBalanceLogic {
	return &GetBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBalanceLogic) GetBalance(in *playground.GetBalanceRequest) (*playground.GetBalanceResponse, error) {
	// todo: add your logic here and delete this line

	return &playground.GetBalanceResponse{}, nil
}
