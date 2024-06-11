package dal

import (
	"database/sql"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"reverse-turing/app/internal/types"
	"reverse-turing/common/errno"
	"reverse-turing/common/model"
	"time"
)

type GameStatus string

const (
	GameStatusCreated GameStatus = "started"
	GameStatusEnded   GameStatus = "ended"
)

type Game struct {
	model.Model
	Code           string           `json:"code" gorm:"column:code;size:32;"`
	DeviceId       string           `json:"deviceId" gorm:"column:device_id;size:32;"`
	Status         GameStatus       `json:"status" gorm:"column:status;size:32;"`
	Progress       string           `json:"progress" gorm:"column:progress;"`
	AiNum          int              `json:"aiNum" gorm:"column:ai_num;type:int;"`
	CharacterNames model.ListString `json:"characterNames" gorm:"column:character_names;type:json;"`
	WhoIsHuman     string           `json:"whoIsHuman" gorm:"column:who_is_human;size:32;"`
	GameStartTime  time.Time        `json:"gameStartTime" gorm:"column:game_start_time;default:null"`
	GameEndTime    sql.NullTime     `json:"gameEndTime" gorm:"column:game_end_time;default:null"`

	AiPlayers []*GameCharacter `json:"aiPlayers" gorm:"-"`
	Player    *GameCharacter   `json:"player" gorm:"-"`
}

func (c *Game) Convert() interface{} {
	var pbResp types.GameResp
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

	if c.AiPlayers != nil {
		pbResp.AiPlayers = make([]types.GameCharacter, len(c.AiPlayers))
		for i, v := range c.AiPlayers {
			pbResp.AiPlayers[i] = *v.Convert().(*types.GameCharacter)
		}
	}
	if c.Player != nil {
		pbResp.Player = *c.Player.Convert().(*types.GameCharacter)
	}

	return &pbResp
}
