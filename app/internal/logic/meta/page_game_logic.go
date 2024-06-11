package meta

import (
	"context"

	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 玩家的游戏列表
func NewPageGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageGameLogic {
	return &PageGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PageGameLogic) PageGame(req *types.EmptyReq) (resp *types.GameAllResp, err error) {
	// todo: add your logic here and delete this line

	return
}
