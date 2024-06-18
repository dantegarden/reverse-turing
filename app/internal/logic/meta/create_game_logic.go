package meta

import (
	"context"
	"github.com/hu-1996/gormx"
	"github.com/samber/lo"
	"gorm.io/gorm"
	"math/rand/v2"
	"reverse-turing/app/internal/dal"
	"reverse-turing/common/errno"
	"reverse-turing/common/utils"
	"time"

	"reverse-turing/app/internal/svc"
	"reverse-turing/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建游戏
func NewCreateGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGameLogic {
	return &CreateGameLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGameLogic) CreateGame(req *types.GameCreateReq) (resp *types.GameResp, err error) {
	deviceId, err := utils.GetDeviceInfo(l.ctx)
	if err != nil {
		return nil, err
	}

	aiNum := req.AiNum
	// 判断aiNum是否是大于2的偶数
	if aiNum < 2 || aiNum%2 != 0 {
		return nil, errno.InvalidAiNumErr
	}

	var newGame *dal.Game

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		newGame = &dal.Game{
			DeviceId:      deviceId,
			Code:          lo.RandomString(6, lo.NumbersCharset),
			Status:        dal.GameStatusCreated,
			GameStartTime: time.Now(),
			Progress:      "0",
			AiNum:         aiNum,
		}

		_, err = gormx.TxInsert[dal.Game](tx, newGame)
		if err != nil {
			return err
		}

		// 查找characters
		var characters []*dal.Character
		err = tx.Raw("SELECT * FROM characters WHERE enable = ? and deleted_at is null ORDER BY RAND() LIMIT ?", true, aiNum+1).Scan(&characters).Error
		if err != nil {
			return err
		}
		if len(characters) < 1 {
			return errno.NoAvailableCharacterErr
		}

		// 分组
		aiCharacters := characters[:aiNum]
		playerCharacter := characters[aiNum]

		// 查找agents
		var agents []*dal.Agent
		err = tx.Raw("SELECT * FROM agents WHERE enable = ? and deleted_at is null ORDER BY RAND() LIMIT ?", true, aiNum+1).Scan(&agents).Error
		if err != nil {
			return err
		}
		if len(agents) < 1 {
			return errno.NoAvailableAgentErr
		}

		var aiPlayers []*dal.GameCharacter
		// 将characters和agents随机组合，分配给aiPlayers
		for _, character := range aiCharacters {
			randIndex := rand.IntN(len(agents))
			agent := agents[randIndex]
			aiCharacter := &dal.GameCharacter{
				GameId:      newGame.ID,
				CharacterId: character.ID,
				AgentId:     agent.ID,
				IsAi:        true,
				Character:   character,
				Agent:       agent,
			}
			aiPlayers = append(aiPlayers, aiCharacter)
		}

		playerAgent := agents[rand.IntN(len(agents))]
		player := &dal.GameCharacter{
			GameId:      newGame.ID,
			CharacterId: playerCharacter.ID,
			Character:   playerCharacter,
			IsAi:        false,
			AgentId:     playerAgent.ID,
			Agent:       playerAgent,
		}

		var allGameCharacters []*dal.GameCharacter
		allGameCharacters = append(allGameCharacters, aiPlayers...)
		allGameCharacters = append(allGameCharacters, player)
		_, err = gormx.TxInsertBatches[dal.GameCharacter](tx, allGameCharacters)
		if err != nil {
			return err
		}

		newGame.CharacterNames = lo.Map[*dal.GameCharacter, string](allGameCharacters, func(_gameCharacter *dal.GameCharacter, _ int) string {
			return _gameCharacter.Character.Name
		})

		_, err = gormx.TxUpdate[dal.Game](tx, newGame)
		if err != nil {
			return err
		}

		newGame.AiPlayers = aiPlayers
		newGame.Player = player

		return nil
	})

	if err != nil {
		return nil, err
	}

	resp = newGame.Convert().(*types.GameResp)
	return
}
