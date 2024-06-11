package dal

import (
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"reverse-turing/app/internal/types"
	"reverse-turing/common/errno"
	"reverse-turing/common/model"
	"time"
)

type GameCharacter struct {
	model.Model
	IsAi        bool   `json:"isAi" gorm:"column:is_ai;"`
	CharacterId string `json:"characterId" gorm:"column:character_id;"`
	AgentId     string `json:"agentId" gorm:"column:agent_id;"`
	GameId      string `json:"gameId" gorm:"column:game_id;"`
	VoteTo      string `json:"voteTo" gorm:"column:vote_to;"`
	FinalVotes  int    `json:"finalVotes" gorm:"column:final_votes;"`

	Character *Character `json:"character" gorm:"-"`
	Agent     *Agent     `json:"agent" gorm:"-"`
}

func (c *GameCharacter) Convert() interface{} {
	var pbResp types.GameCharacter
	_ = copier.CopyWithOption(&pbResp, c, copier.Option{Converters: []copier.TypeConverter{
		{
			SrcType: time.Time{},
			DstType: copier.String,
			Fn: func(src interface{}) (interface{}, error) {
				s, ok := src.(time.Time)

				if !ok {
					logx.Errorf("src type not matching")
					return nil, errno.ServerErr
				}

				return s.Format(time.DateTime), nil
			},
		}},
		DeepCopy: true,
	})

	if c.Character != nil {
		pbResp.Character = c.Character.Convert().(*types.Character)
	}
	if c.Agent != nil {
		pbResp.Agent = c.Agent.Convert().(*types.Agent)
	}
	return &pbResp
}
