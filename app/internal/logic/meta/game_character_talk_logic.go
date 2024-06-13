package meta

import (
	"context"
	"errors"
	"fmt"
	"github.com/hu-1996/gormx"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhouyangtingwen/dify-sdk-go"
	"gorm.io/gorm"
	"net/http"
	"reverse-turing/app/internal/dal"
	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"
	"reverse-turing/common/errno"
	"reverse-turing/common/utils"
	"strconv"
	"strings"
)

type GameCharacterTalkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 与AI角色对话
func NewGameCharacterTalkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GameCharacterTalkLogic {
	return &GameCharacterTalkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GameCharacterTalkLogic) GameCharacterTalk(req *types.GameCharacterTalkReq, w http.ResponseWriter) (resp *types.EmptyResp, err error) {
	deviceId, err := utils.GetDeviceInfo(l.ctx)
	if err != nil {
		return nil, err
	}

	game, err := gormx.SelectById[dal.Game](req.GameId)
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

	// 准备调dify
	difyConf := utils.GetDifyConfig(l.ctx, gameCharacter.Agent.Endpoint, gameCharacter.Agent.ApiKey)
	// 构建请求体
	characterNames := utils.RemoveStringCreateNew(game.CharacterNames, gameCharacter.Character.Name)
	query, err := dal.GetQuery(dal.TalkType(req.TalkType), req.Params...)
	if err != nil {
		return nil, err
	}
	difyReq := &dify.ChatMessageRequest{
		Inputs: map[string]interface{}{
			"character":        gameCharacter.Character.Name,
			"ai_num":           strconv.Itoa(game.AiNum),
			"other_characters": strings.Join(characterNames, "、"),
			"positioning":      gameCharacter.Character.Positioning,
		},
		ResponseMode: "streaming",
		Query:        query,
		User:         deviceId,
	}

	err = utils.SendDifyMessage(l.ctx, difyConf, difyReq, func(d string) {
		_, err := fmt.Fprintf(w, "%s", d)
		if err != nil {
			return
		}
		w.(http.Flusher).Flush()
	})
	if err != nil {
		return nil, err
	}

	return
}
