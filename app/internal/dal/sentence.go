package dal

import (
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"reverse-turing/app/internal/types"
	"reverse-turing/common/errno"
	"reverse-turing/common/model"
	"time"
)

type Sentence struct {
	model.Model
	GameCharacterId string `json:"gameCharacterId" gorm:"column:game_character_id;size:32;"`
	CharacterId     string `json:"characterId" gorm:"column:character_id;size:32;"`
	AgentId         string `json:"agentId" gorm:"column:agent_id;size:32;"`
	Name            string `json:"name" gorm:"column:name;size:32;"`
	Content         string `json:"content" gorm:"column:content;size:255;"`
	GameId          string `json:"gameId" gorm:"column:game_id;size:32;index"`
	TalkType        string `json:"talkType" gorm:"column:talk_type;size:100;"`
}

func (c *Sentence) Convert() interface{} {
	var pbResp types.SentenceResp
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

	return &pbResp
}
