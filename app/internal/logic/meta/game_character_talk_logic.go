package meta

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hu-1996/gormx"
	"github.com/samber/lo"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zhouyangtingwen/dify-sdk-go"
	"gorm.io/gorm"
	"net/http"
	"reverse-turing/app/internal/dal"
	"reverse-turing/app/internal/sse"
	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"
	"reverse-turing/common/consts"
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

func (l *GameCharacterTalkLogic) GameCharacterTalk(req *types.GameCharacterTalkReq, r *http.Request, w http.ResponseWriter) (resp *types.EmptyResp, err error) {
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
	difyConf := GetDifyConfig(l.ctx, gameCharacter.Agent.Endpoint, gameCharacter.Agent.ApiKey)
	// 构建请求体
	characterNames := utils.RemoveStringCreateNew(game.CharacterNames, gameCharacter.Character.Name)
	talkType := dal.TalkType(req.TalkType)
	query, err := dal.GetQuery(talkType, req.Params...)
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

	recvContent, err := SendDifyStreamMessage(l.ctx, difyConf, difyReq, w, GetEndMark(talkType), true)
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
		Content:         recvContent,
		TalkType:        req.TalkType,
	}
	_, err = gormx.Insert[dal.Sentence](sentence)
	if err != nil {
		return nil, err
	}

	// 是否需要裁决
	if lo.Contains([]dal.TalkType{dal.TalkTypeAsk, dal.TalkTypeVote}, talkType) {

		var judgeDifyConf *dify.Client
		var judgeDifyReq *dify.ChatMessageRequest

		switch talkType {
		case dal.TalkTypeAsk:
			judgeDifyConf = GetDifyConfig(l.ctx, l.svcCtx.Config.AiBot.AskJudge.Endpoint, l.svcCtx.Config.AiBot.AskJudge.ApiKey)
			judgeDifyReq = &dify.ChatMessageRequest{
				Inputs: map[string]interface{}{
					"characters": strings.Join(characterNames, "、"),
				},
				ResponseMode: "streaming",
				Query:        fmt.Sprintf("%s：%s", gameCharacter.Character.Name, recvContent),
				User:         deviceId,
			}

		case dal.TalkTypeVote:
			judgeDifyConf = GetDifyConfig(l.ctx, l.svcCtx.Config.AiBot.VoteJudge.Endpoint, l.svcCtx.Config.AiBot.VoteJudge.ApiKey)
			judgeDifyReq = &dify.ChatMessageRequest{
				Inputs:       map[string]interface{}{},
				ResponseMode: "streaming",
				Query:        fmt.Sprintf("%s：%s", gameCharacter.Character.Name, recvContent),
				User:         deviceId,
			}
		}

		judgeRecvContent, err := SendDifyStreamMessage(l.ctx, judgeDifyConf, judgeDifyReq, w, consts.EndMarkDone, false)
		if err != nil {
			return nil, err
		}
		logx.Infof("judgeRecvContent: %s", judgeRecvContent)

		var judgeResp types.GameJudgeResp
		err = DeserializeJSON(judgeRecvContent, &judgeResp)
		if err != nil {
			return nil, err
		}

		judgeResp.JudgeType = req.TalkType
		marshal, _ := json.Marshal(judgeResp)
		Publish(w, &sse.Event{
			Data: marshal,
		})
		Publish(w, &sse.Event{Data: []byte(consts.EndMarkDone)})
	}
	return &types.EmptyResp{}, nil
}

func GetEndMark(taskType dal.TalkType) string {
	switch taskType {
	case dal.TalkTypeOpening, dal.TalkTypeAnswer:
		return consts.EndMarkDone
	case dal.TalkTypeAsk:
		return consts.EndMarkAsk
	case dal.TalkTypeVote:
		return consts.EndMarkVote
	}
	return consts.EndMarkDone
}

func DeserializeJSON(input string, judgeResp *types.GameJudgeResp) error {
	// 处理 markdown 格式的 json 字符串
	if strings.HasPrefix(input, "```json") {
		input = strings.TrimPrefix(input, "```json")
		input = strings.TrimSuffix(input, "```")
		input = strings.TrimSpace(input)
	}

	err := json.Unmarshal([]byte(input), judgeResp)
	return err
}
