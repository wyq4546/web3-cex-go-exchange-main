package logic

import (
	"context"

	"solana/internal/svc"
	"grpc-common/solana/types/playground"

	"github.com/zeromicro/go-zero/core/logx"
)

type AirdropLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAirdropLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AirdropLogic {
	return &AirdropLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AirdropLogic) Airdrop(in *playground.AirdropRequest) (*playground.AirdropResponse, error) {
	// todo: add your logic here and delete this line

	return &playground.AirdropResponse{}, nil
}
