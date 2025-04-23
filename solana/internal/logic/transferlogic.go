package logic

import (
	"context"

	"solana/internal/svc"
	"grpc-common/solana/types/playground"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferLogic {
	return &TransferLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransferLogic) Transfer(in *playground.TransferRequest) (*playground.TransferResponse, error) {
	// todo: add your logic here and delete this line

	return &playground.TransferResponse{}, nil
}
