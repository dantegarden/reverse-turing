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

type PlayerTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 用户对话
func NewPlayerTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlayerTalkLogic {
	return &PlayerTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PlayerTalkLogic) PlayerTalk(req *types.GamePlayerTalkReq) (resp *types.EmptyResp, err error) {
	_, err = utils.GetDeviceInfo(l.ctx)
	if err != nil {
		return nil, err
	}

	_, err = gormx.SelectById[dal.Game](req.GameId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.GameNotFoundErr
		}
		return nil, err
	}

	gameCharacter, err := gormx.SelectById[dal.GameCharacter](req.GameCharacterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.GameCharacterNotFoundErr
		}
		return nil, err
	}
	if gameCharacter.GameId != req.GameId {
		return nil, errno.GameCharacterNotFoundErr
	}

	err = CompleteGameCharacter(gameCharacter)
	if err != nil {
		return nil, err
	}

	// 记录句子
	sentence := &dal.Sentence{
		GameId:          req.GameId,
		GameCharacterId: gameCharacter.ID,
		CharacterId:     gameCharacter.CharacterId,
		AgentId:         gameCharacter.AgentId,
		Name:            gameCharacter.Character.Name,
		Content:         req.Content,
		TalkType:        req.TalkType,
	}
	_, err = gormx.Insert[dal.Sentence](sentence)
	if err != nil {
		return nil, err
	}

	return
}
