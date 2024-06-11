package dal

import (
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"reverse-turing/app/internal/types"
	"reverse-turing/common/errno"
	"reverse-turing/common/model"
	"time"
)

type Character struct {
	model.Model
	Name        string `json:"name" gorm:"column:name;size:32;"`
	Positioning string `json:"positioning" gorm:"column:positioning;size:32;"`
	Avatar      string `json:"avatar" gorm:"column:avatar;"`
	Portrait    string `json:"portrait" gorm:"column:portrait;"`
	Profile     string `json:"profile" gorm:"column:profile;"`
	Enable      bool   `json:"enable" gorm:"column:enable;"`
}

func (c *Character) Convert() interface{} {
	var pbResp types.Character
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
