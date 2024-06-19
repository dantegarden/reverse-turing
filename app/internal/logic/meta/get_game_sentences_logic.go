package meta

import (
	"context"
	"errors"
	"github.com/hu-1996/gormx"
	"gorm.io/gorm"
	"reverse-turing/app/internal/dal"
	"reverse-turing/common/errno"
	"reverse-turing/common/utils"

	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameSentencesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取游戏对话记录
func NewGetGameSentencesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameSentencesLogic {
	return &GetGameSentencesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGameSentencesLogic) GetGameSentences(req *types.GamePathReq) (resp []types.SentenceResp, err error) {
	_, err = utils.GetDeviceInfo(l.ctx)
	if err != nil {
		return nil, err
	}

	game, err := gormx.SelectById[dal.Game](req.Id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.GameNotFoundErr
		}
		return nil, err
	}

	sentences, err := gormx.SelectListConvert[dal.Sentence, types.SentenceResp]("id asc", "game_id = ?", game.ID)
	if err != nil {
		return nil, err
	}
	resp = make([]types.SentenceResp, len(sentences))
	for i, ptr := range sentences {
		if ptr != nil { // 确保指针不为空
			resp[i] = *ptr
		}
	}

	return
}
