package meta

import (
	"context"
	"errors"
	"github.com/hu-1996/gormx"
	"gorm.io/gorm"
	"reverse-turing/app/internal/dal"
	"reverse-turing/common/errno"

	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 玩家的游戏列表
func NewGetGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameLogic {
	return &GetGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGameLogic) GetGame(req *types.GamePathReq) (resp *types.GameResp, err error) {
	id := req.Id
	game, err := gormx.SelectById[dal.Game](id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = errors.New("游戏不存在")
		return
	}

	var aiPlayers []*dal.GameCharacter
	gcs, err := gormx.SelectList[dal.GameCharacter]("", "game_id = ?", game.ID)
	for _, gc := range gcs {
		err2 := CompleteGameCharacter(gc)
		if err2 != nil {
			err = err2
			return
		}
		if gc.IsAi {
			aiPlayers = append(aiPlayers, gc)
		} else {
			game.Player = gc
		}
	}
	game.AiPlayers = aiPlayers

	resp = game.Convert().(*types.GameResp)
	return
}

func CompleteGameCharacter(gc *dal.GameCharacter) error {
	character, err := gormx.SelectById[dal.Character](gc.CharacterId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.CharacterNotFoundErr
		}
		return err
	}
	agent, err := gormx.SelectById[dal.Agent](gc.AgentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errno.AgentNotFoundErr
		}
		return err
	}

	gc.Character = character
	gc.Agent = agent
	return nil
}
